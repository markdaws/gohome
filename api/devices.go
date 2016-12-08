package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/go-home-iot/connection-pool"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/feature"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
	errExt "github.com/pkg/errors"
)

// RegisterDeviceHandlers registers the REST API routes relating to devices
func RegisterDeviceHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/devices",
		apiDevicesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/devices",
		apiAddDeviceHandler(s.systemSavePath, s.system)).Methods("POST")
	r.HandleFunc("/api/v1/devices/{id}",
		apiDeviceHandlerDelete(s.systemSavePath, s.system)).Methods("DELETE")
	r.HandleFunc("/api/v1/devices/{id}",
		apiDeviceHandlerUpdate(s.systemSavePath, s.system)).Methods("PUT")
	r.HandleFunc("/api/v1/devices/{id}/features/{fid}",
		apiDeviceUpdateFeatureHandler(s.systemSavePath, s.system)).Methods("PUT")
	r.HandleFunc("/api/v1/devices/{id}/features/{fid}/apply",
		apiDeviceApplyFeaturesAttrsHandler(s.systemSavePath, s.system)).Methods("PUT")
}

func DevicesToJSON(devs map[string]*gohome.Device) []jsonDevice {
	devices := make(devices, len(devs))
	var i int32
	for _, device := range devs {
		var connPoolJSON *jsonConnPool
		if device.Connections != nil {
			connCfg := device.Connections.Config
			connPoolJSON = &jsonConnPool{
				Name:     connCfg.Name,
				PoolSize: int32(connCfg.Size),
			}
		}

		var authJSON *jsonAuth
		if device.Auth != nil {
			authJSON = &jsonAuth{
				Login:    device.Auth.Login,
				Password: device.Auth.Password,
				Token:    device.Auth.Token,
			}
		}

		hubID := ""
		if device.Hub != nil {
			hubID = device.Hub.ID
		}

		var deviceIDs []string
		for _, dev := range device.Devices {
			deviceIDs = append(deviceIDs, dev.ID)
		}

		devices[i] = jsonDevice{
			ID:              device.ID,
			Address:         device.Address,
			Name:            device.Name,
			Description:     device.Description,
			ModelNumber:     device.ModelNumber,
			ModelName:       device.ModelName,
			SoftwareVersion: device.SoftwareVersion,
			ConnPool:        connPoolJSON,
			Auth:            authJSON,
			Type:            string(device.Type),
			HubID:           hubID,
			DeviceIDs:       deviceIDs,
			Features:        device.Features,
		}
		i++
	}
	sort.Sort(devices)
	return devices
}

func apiDevicesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(w).Encode(DevicesToJSON(system.Devices())); err != nil {
			respErr(err, w)
		}
	}
}

func apiDeviceHandlerDelete(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		deviceID := mux.Vars(r)["id"]
		device := system.DeviceByID(deviceID)
		if device == nil {
			respBadRequest(fmt.Sprintf("invalid device ID: %s", deviceID), w)
			return
		}
		system.DeleteDevice(device)
		err := store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "failed to save changes to disk"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiDeviceApplyFeaturesAttrsHandler(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			respBadRequest(fmt.Sprintf("failed to read request body: %s", err), w)
			return
		}

		// Get the attributes and make commands
		var data = make(map[string]*attr.Attribute)

		if err = json.Unmarshal(body, &data); err != nil {
			respBadRequest(fmt.Sprintf("invalid request body: %s", err), w)
			return
		}

		// When deserializing from JSON, the int32 and float32 types are converted
		// to float64 so need to massage them back
		attr.FixJSON(data)

		deviceID := mux.Vars(r)["id"]
		dev := system.DeviceByID(deviceID)
		if dev == nil {
			respBadRequest(fmt.Sprintf("invalid device ID: %s", deviceID), w)
			return
		}

		featureID := mux.Vars(r)["fid"]
		f := system.FeatureByID(featureID)
		if f == nil {
			respBadRequest(fmt.Sprintf("invalid feature ID: %s", featureID), w)
			return
		}

		// Verify that each attribute passed in is valid. The API only cares that you pass in
		// localID and value, the other fields for the attribute are pulled from the feature
		finalAttrs := make(map[string]*attr.Attribute)
		for localID, attribute := range data {
			blankAttr, ok := f.Attrs[localID]
			if !ok {
				respBadRequest(fmt.Sprintf("invalid localID: %s", localID), w)
				return
			}

			finalAttrs[localID] = blankAttr.Clone()
			finalAttrs[localID].Value = attribute.Value
		}

		desc := "FeatureSetAttrs"
		err = system.Services.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.FeatureSetAttrs{
			FeatureID:   featureID,
			FeatureName: f.Name,
			Attrs:       finalAttrs,
		}))

		if err != nil {
			respErr(errExt.Wrap(err, "failed to enqueue FeatureSetAttrs command"), w)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiDeviceUpdateFeatureHandler(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			respBadRequest(fmt.Sprintf("failed to read request body: %s", err), w)
			return
		}

		var data feature.Feature
		if err = json.Unmarshal(body, &data); err != nil {
			respBadRequest(fmt.Sprintf("invalid request body: %s", err), w)
			return
		}

		deviceID := mux.Vars(r)["id"]
		dev := system.DeviceByID(deviceID)
		if dev == nil {
			respBadRequest(fmt.Sprintf("invalid device ID: %s", deviceID), w)
			return
		}

		featureID := mux.Vars(r)["fid"]
		f := system.FeatureByID(featureID)
		if f == nil {
			respBadRequest(fmt.Sprintf("invalid feature ID: %s", featureID), w)
			return
		}

		updatedFeature := feature.Feature{
			ID:          data.ID,
			Type:        data.Type,
			Name:        data.Name,
			Address:     data.Address,
			Description: data.Description,
			DeviceID:    data.DeviceID,
		}

		valErrs := updatedFeature.Validate()
		if valErrs != nil {
			respValErr(&data, data.ID, valErrs, w)
			return
		}

		if data.Type != f.Type {
			//TODO: Need to update commands since they store the type and attrs for the feature which are now all invalid
			fmt.Println("TODO: If you change the type of a feature, then need to update all the commands")

			// Changing the type of the feature, we need to update the attributes
			// associated with the feature
			newFeature := feature.NewFromType(data.ID, data.Type)
			if newFeature == nil {
				valErrs = &validation.Errors{}
				valErrs.Add("unsupported. Can't change to target type", "Type")
				respValErr(&data, data.ID, valErrs, w)
				return
			}

			f.Attrs = newFeature.Attrs
		}

		// Validation passed now we can update the actual object
		f.Type = data.Type
		f.Name = data.Name
		f.Address = data.Address
		f.Description = data.Description

		err = store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "error writing changes to disk"), w)
			return
		}

		// Don't support chaning attrs at this moment
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(f)
	}
}

func apiAddDeviceHandler(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			respBadRequest(fmt.Sprintf("failed to read request body: %s", err), w)
			return
		}

		var data jsonDevice
		if err = json.Unmarshal(body, &data); err != nil {
			respBadRequest(fmt.Sprintf("invalid request body: %s", err), w)
			return
		}

		if dev := system.DeviceByID(data.ID); dev == nil {
			respBadRequest(fmt.Sprintf("trying to add a device with a duplicate ID: %s", data.ID), w)
			return
		}

		var auth *gohome.Auth
		if data.Auth != nil {
			auth = &gohome.Auth{
				Login:    data.Auth.Login,
				Password: data.Auth.Password,
				Token:    data.Auth.Token,
			}
		}

		var hub *gohome.Device
		if data.HubID != "" {
			if hub = system.DeviceByID(data.HubID); hub == nil {
				respBadRequest(fmt.Sprintf("invalid hub ID: %s", data.HubID), w)
				return
			}
		}

		d := gohome.NewDevice(
			data.ID,
			data.Name,
			data.Description,
			data.ModelNumber,
			data.ModelName,
			data.SoftwareVersion,
			data.Address,
			hub,
			nil,
			nil,
			auth,
		)
		d.Type = gohome.DeviceType(data.Type)

		var connPoolCfg *pool.Config
		if data.ConnPool != nil {
			connPoolCfg = &pool.Config{
				Name: data.ConnPool.Name,
				Size: int(data.ConnPool.PoolSize),

				//TODO Serialize
				RetryDuration: time.Second * 10,
			}

			network := system.Extensions.FindNetwork(system, d)
			if network != nil {
				connFactory, err := network.NewConnection(system, d)
				if err == nil {
					connPoolCfg.NewConnection = connFactory
				}
			}
			d.Connections = pool.NewPool(*connPoolCfg)
		}

		cmdBuilder := system.Extensions.FindCmdBuilder(system, d)
		d.CmdBuilder = cmdBuilder

		valErrs := d.Validate()
		if valErrs != nil {
			respValErr(&data, data.ID, valErrs, w)
			return
		}

		for _, f := range data.Features {
			if ft := system.FeatureByID(f.ID); ft == nil {
				respBadRequest(fmt.Sprintf("trying to add duplicate feature %s", f.ID), w)
				return
			}

			valErrs := f.Validate()
			if valErrs != nil {
				respValErr(&f, f.ID, valErrs, w)
				return
			}

			// Have to fix JSON unmarshal, converting types to float64
			attr.FixJSON(f.Attrs)

			err = d.AddFeature(f)
			if err != nil {
				respBadRequest(fmt.Sprintf("error adding feature to device %s", err), w)
				return
			}
		}

		// Now we have created everything and know all the validation passes we can
		// commit all the changes at once
		system.AddDevice(d)
		for _, feature := range data.Features {
			system.AddFeature(feature)
		}

		err = system.InitDevice(d)
		if err != nil {
			log.E("Failed to init device on add: %s", err)
		}

		err = store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "error writing changes to disk"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

func apiDeviceHandlerUpdate(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			respBadRequest(fmt.Sprintf("failed to read request body: %s", err), w)
			return
		}

		var data jsonDevice
		if err = json.Unmarshal(body, &data); err != nil {
			respBadRequest(fmt.Sprintf("failed to parse JSON body: %s", err), w)
			return
		}

		d := system.DeviceByID(data.ID)
		if d == nil {
			respBadRequest(fmt.Sprintf("invalid device ID: %s", data.ID), w)
			return
		}

		updatedDev := gohome.NewDevice(
			data.ID,
			data.Name,
			data.Description,
			data.ModelNumber,
			"",
			"",
			data.Address,
			nil,
			nil,
			nil,
			nil)

		errors := updatedDev.Validate()
		if errors != nil {
			respValErr(&data, data.ID, errors, w)
			return
		}

		// Validated, set the fields
		d.Name = data.Name
		d.Description = data.Description
		d.ModelNumber = data.ModelNumber

		addressChanged := d.Address != data.Address
		d.Address = data.Address
		d.Type = gohome.DeviceType(data.Type)

		err = store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "failed to save new settings to disk"), w)
			return
		}

		// If the address changed then we need to stop all services associated
		// with the device and start them again using the new address
		if addressChanged {
			//TODO: Finish this, pattern for stopping devices
			system.StopDevice(d)
			system.InitDevice(d)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		jsonDevices := DevicesToJSON(map[string]*gohome.Device{d.ID: d})
		json.NewEncoder(w).Encode(jsonDevices[0])
	}
}

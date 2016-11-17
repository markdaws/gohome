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
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
	errExt "github.com/pkg/errors"
)

// RegisterDeviceHandlers registers the REST API routes relating to devices
func RegisterDeviceHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/devices",
		apiDevicesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/devices",
		apiAddDeviceHandler(s.systemSavePath, s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/devices/{id}",
		apiDeviceHandlerDelete(s.systemSavePath, s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/devices/{id}",
		apiDeviceHandlerUpdate(s.systemSavePath, s.system, s.recipeManager)).Methods("PUT")
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

		jsonButtons := ButtonsToJSON(device.Buttons)
		jsonZones := ZonesToJSON(device.Zones)
		jsonSensors := SensorsToJSON(device.Sensors)

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
			Zones:           jsonZones,
			Sensors:         jsonSensors,
			Buttons:         jsonButtons,
			ConnPool:        connPoolJSON,
			Auth:            authJSON,
			Type:            string(device.Type),
			HubID:           hubID,
			DeviceIDs:       deviceIDs,
		}
		i++
	}
	sort.Sort(devices)
	return devices
}

func apiDevicesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(w).Encode(DevicesToJSON(system.Devices)); err != nil {
			respErr(err, w)
		}
	}
}

func apiDeviceHandlerDelete(
	savePath string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		deviceID := mux.Vars(r)["id"]
		device, ok := system.Devices[deviceID]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid device ID: %s", deviceID), w)
			return
		}
		system.DeleteDevice(device)
		err := store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			respErr(errExt.Wrap(err, "failed to save changes to disk"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiAddDeviceHandler(
	savePath string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
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

		if _, ok := system.Devices[data.ID]; ok {
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
		var foundHub bool
		if data.HubID != "" {
			if hub, foundHub = system.Devices[data.HubID]; !foundHub {
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

		newButtons := make([]*gohome.Button, 0, len(data.Buttons))
		for _, button := range data.Buttons {
			btn := &gohome.Button{
				Address:     button.Address,
				ID:          button.ID,
				Name:        button.Name,
				Description: button.Description,
				Device:      d,
			}
			err := d.AddButton(btn)
			if err != nil {
				respBadRequest(fmt.Sprintf("failed to add button to device: %s", err), w)
				return
			}
			newButtons = append(newButtons, btn)
		}

		newZones := make([]*zone.Zone, 0, len(data.Zones))
		for _, zn := range data.Zones {
			if _, ok := system.Zones[zn.ID]; ok {
				respBadRequest(fmt.Sprintf("trying to add duplicate zone %s", zn.ID), w)
				return
			}

			z := &zone.Zone{
				ID:          zn.ID,
				Address:     zn.Address,
				Name:        zn.Name,
				Description: zn.Description,
				DeviceID:    zn.DeviceID,
				Type:        zone.TypeFromString(zn.Type),
				Output:      zone.OutputFromString(zn.Output),
			}

			valErrs := z.Validate()
			if valErrs != nil {
				respValErr(&zn, zn.ID, valErrs, w)
				return
			}

			errors := d.AddZone(z)
			if errors != nil {
				if valErrs, ok := errors.(*validation.Errors); ok {
					respValErr(&zn, zn.ID, valErrs, w)
				} else {
					respBadRequest(fmt.Sprintf("error adding zone to device %s", errors), w)
				}
				return
			}
			newZones = append(newZones, z)
		}

		newSensors := make([]*gohome.Sensor, 0, len(data.Sensors))
		for _, sen := range data.Sensors {
			if _, ok := system.Sensors[sen.ID]; ok {
				respBadRequest(fmt.Sprintf("trying to add duplicate sensor %s", sen.ID), w)
				return
			}

			sensor := &gohome.Sensor{
				ID:          sen.ID,
				Name:        sen.Name,
				Description: sen.Description,
				Address:     sen.Address,
				DeviceID:    sen.DeviceID,
				Attr: gohome.SensorAttr{
					Name:     sen.Attr.Name,
					DataType: gohome.SensorDataType(sen.Attr.DataType),
					States:   sen.Attr.States,
				},
			}

			valErrs := sensor.Validate()
			if valErrs != nil {
				respValErr(&sen, sen.ID, valErrs, w)
				return
			}

			err = d.AddSensor(sensor)
			if err != nil {
				respBadRequest(fmt.Sprintf("error adding sensor to device %s", err), w)
				return
			}
			newSensors = append(newSensors, sensor)
		}

		// Now we have created everything and know all the validation passes we can
		// commit all the changes at once
		system.AddDevice(d)
		for _, button := range newButtons {
			system.AddButton(button)
		}
		for _, zone := range newZones {
			system.AddZone(zone)
		}
		for _, sensor := range newSensors {
			system.AddSensor(sensor)
		}

		err = system.InitDevice(d)
		if err != nil {
			log.E("Failed to init device on add: %s", err)
		}

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			respErr(errExt.Wrap(err, "error writing changes to disk"), w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

func apiDeviceHandlerUpdate(
	savePath string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
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

		d, ok := system.Devices[data.ID]
		if !ok {
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

		err = store.SaveSystem(savePath, system, recipeManager)
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
		json.NewEncoder(w).Encode(data)
	}
}

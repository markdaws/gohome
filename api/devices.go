package api

import (
	"encoding/json"
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
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		system.DeleteDevice(device)
		err := store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			log.V("Failed to ready request body %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data jsonDevice
		if err = json.Unmarshal(body, &data); err != nil {
			log.V("error unmarhsaling device %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, ok := system.Devices[data.ID]; ok {
			log.V("trying to add duplicate device %s", data.ID)
			w.WriteHeader(http.StatusBadRequest)
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

		d := gohome.NewDevice(
			data.ModelNumber,
			data.ModelName,
			data.SoftwareVersion,
			data.Address,
			data.ID,
			data.Name,
			data.Description,
			nil,
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
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(validation.NewErrorJSON(&data, data.ID, valErrs))
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
				log.V("failed to add button to device: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			newButtons = append(newButtons, btn)
		}

		newZones := make([]*zone.Zone, 0, len(data.Zones))
		for _, zn := range data.Zones {
			if _, ok := system.Zones[zn.ID]; ok {
				log.V("trying to add duplicate zone %s", zn.ID)
				w.WriteHeader(http.StatusBadRequest)
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
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&zn, zn.ID, valErrs))
				return
			}

			errors := d.AddZone(z)
			if errors != nil {
				if valErrs, ok := errors.(*validation.Errors); ok {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&zn, zn.ID, valErrs))
				} else {
					log.V("error adding zone to device %s", errors)
					w.WriteHeader(http.StatusBadRequest)
				}
				return
			}
			newZones = append(newZones, z)
		}

		newSensors := make([]*gohome.Sensor, 0, len(data.Sensors))
		for _, sen := range data.Sensors {
			if _, ok := system.Sensors[sen.ID]; ok {
				log.V("trying to add duplicate sensor %s", sen.ID)
				w.WriteHeader(http.StatusBadRequest)
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
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&sen, sen.ID, valErrs))
				return
			}

			err = d.AddSensor(sensor)
			if err != nil {
				log.V("error adding sensor to device %s", err)
				w.WriteHeader(http.StatusBadRequest)
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
			log.V("error writing changes to disk %s", err)
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data jsonDevice
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		d, ok := system.Devices[data.ID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedDev := gohome.NewDevice(
			data.ModelNumber,
			"",
			"",
			data.Address,
			data.ID,
			data.Name,
			data.Description,
			nil,
			nil,
			nil,
			nil)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		errors := updatedDev.Validate()
		if errors != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(validation.NewErrorJSON(&data, data.ID, errors))
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
			w.WriteHeader(http.StatusInternalServerError)
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

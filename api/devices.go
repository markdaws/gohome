package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/go-home-iot/connection-pool"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
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

func apiDevicesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		devices := make(devices, len(system.Devices), len(system.Devices))
		var i int32
		for _, device := range system.Devices {
			devices[i] = jsonDevice{
				Address:     device.Address,
				ID:          device.ID,
				Name:        device.Name,
				Description: device.Description,
				ModelNumber: device.ModelNumber,
				Type:        string(device.Type),
			}
			i++
		}
		sort.Sort(devices)
		if err := json.NewEncoder(w).Encode(devices); err != nil {
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

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data jsonDevice
		if err = json.Unmarshal(body, &data); err != nil {
			fmt.Printf("%s\n", err)
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
			"",
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

		errors := system.AddDevice(d)
		if errors != nil {
			if valErrs, ok := errors.(*validation.Errors); ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&data, data.ClientID, valErrs))
			} else {
				//Other kind of errors, TODO: log
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		err = system.InitDevice(d)
		if err != nil {
			log.E("Failed to init device on add: %s", err)
		}

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		data.ClientID = ""
		data.ID = d.ID
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

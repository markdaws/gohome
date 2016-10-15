package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/discovery"
)

func apiDiscoveryHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		modelNumber := vars["modelNumber"]

		//TODO: fix, This is blocking
		devs, err := discovery.Devices(system, modelNumber)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonDevices := make(devices, len(devs))
		for i, device := range devs {

			//TODO: This API shouldn't be sending back client ids, the client should
			//make these values up
			jsonZones := make(zones, len(device.Zones))
			var j int
			for _, zn := range device.Zones {
				jsonZones[j] = jsonZone{
					Address:     zn.Address,
					ID:          zn.ID,
					Name:        zn.Name,
					Description: zn.Description,
					DeviceID:    device.ID,
					Type:        zn.Type.ToString(),
					Output:      zn.Output.ToString(),
					ClientID:    modelNumber + "_zone_" + strconv.Itoa(j),
				}
				j++
			}

			var cmdBuilderJson *jsonCmdBuilder
			if device.CmdBuilder != nil {
				cmdBuilderJson = &jsonCmdBuilder{ID: device.CmdBuilder.ID()}
			}
			var connPoolJson *jsonConnPool
			if device.Connections != nil {
				connCfg := device.Connections.Config()
				connPoolJson = &jsonConnPool{
					Name:           connCfg.Name,
					PoolSize:       int32(connCfg.Size),
					ConnectionType: connCfg.ConnectionType,
					TelnetPingCmd:  connCfg.TelnetPingCmd,
					Address:        connCfg.Address,
				}
			}
			jsonDevices[i] = jsonDevice{
				ID:          device.ID,
				Address:     device.Address,
				Name:        device.Name,
				Description: device.Description,
				ModelNumber: device.ModelNumber,
				Token:       "",
				ClientID:    modelNumber + "_" + strconv.Itoa(i),
				Zones:       jsonZones,
				CmdBuilder:  cmdBuilderJson,
				ConnPool:    connPoolJson,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(jsonDevices); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		/*TODO: Remove
		json.NewEncoder(w).Encode(struct {
			Location string `json:"location"`
		}{Location: data["location"]})
		*/
	}
}

/*
func apiDiscoveryZoneHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//This is blocking, waits 5 seconds
		zs, err := discovery.Zones(vars["modelNumber"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonZones := make(zones, len(zs))
		for i, zone := range zs {
			jsonZones[i] = jsonZone{
				Address:     zone.Address,
				Name:        zone.Name,
				Description: zone.Description,
				Type:        zone.Type.ToString(),
				Output:      zone.Output.ToString(),
			}
		}
		sort.Sort(jsonZones)
		if err := json.NewEncoder(w).Encode(jsonZones); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
*/

func apiDiscoveryTokenHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//TODO: Make non-blocking: this is blocking
		token, err := discovery.DiscoverToken(vars["modelNumber"], r.URL.Query().Get("address"))
		if err != nil {
			if err == discovery.ErrUnauthorized {
				// Let the caller know this was a specific kind of error
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(struct {
					Unauthorized bool `json:"unauthorized"`
				}{Unauthorized: true})
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{Token: token})
	}
}

func apiDiscoveryAccessHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//TODO: Make non-blocking: this is blocking
		err := discovery.VerifyConnection(
			vars["modelNumber"],
			r.URL.Query().Get("address"),
			r.URL.Query().Get("token"),
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

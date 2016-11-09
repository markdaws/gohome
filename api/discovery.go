package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterDiscoveryHandlers registers all of the discovery specific API REST routes
func RegisterDiscoveryHandlers(r *mux.Router, s *apiServer) {

	// Get a list of all the devices that we can discover
	r.HandleFunc("/api/v1/discovery/discoverers",
		apiListDiscoveryHandler(s.system)).Methods("GET")

	// Scan the network for all devices corresponding to the discovery ID
	r.HandleFunc("/api/v1/discovery/discoverers/{discovererID}",
		apiDiscoveryHandler(s.system)).Methods("GET")

	r.HandleFunc("/api/v1/discovery/discoverers/{discovererID}",
		apiFromStringDiscoveryHandler(s.system)).Methods("POST")

}

func apiListDiscoveryHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		infos := system.Extensions.ListDiscoverers(system)

		jsonInfos := make([]jsonDiscovererInfo, len(infos), len(infos))
		for i, info := range infos {
			jsonInfos[i] = jsonDiscovererInfo{
				ID:          info.ID,
				Name:        info.Name,
				Description: info.Description,
				Type:        info.Type,
				PreScanInfo: info.PreScanInfo,
			}
		}
		if err := json.NewEncoder(w).Encode(jsonInfos); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func writeDiscoveryResults(result *gohome.DiscoveryResults, w http.ResponseWriter) {
	//TODO: Serializing scenes, if they don't have ids then we need to use fake ids in the commands

	// Need to serialize the scenes, use handy functions from scenes.go
	inputScenes := make(map[string]*gohome.Scene)
	for i, scene := range result.Scenes {
		// don't have ids at this point
		inputScenes[strconv.Itoa(i)] = scene
	}

	jsonScenes := scenesToJSON(inputScenes)
	jsonDevices := make(devices, len(result.Devices))
	for i, device := range result.Devices {

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
				ClientID:    strconv.Itoa(i) + "_zone_" + strconv.Itoa(j),
			}
			j++
		}

		jsonSensors := make(sensors, len(device.Sensors))
		j = 0
		for _, sen := range device.Sensors {
			jsonSensors[j] = jsonSensor{
				ID:          sen.ID,
				Name:        sen.Name,
				Description: sen.Description,
				Address:     sen.Address,

				//TODO: Shouldn't be setting ClientID here
				ClientID: strconv.Itoa(i) + "_sensor_" + strconv.Itoa(j),

				Attr: jsonSensorAttr{
					Name:     sen.Attr.Name,
					DataType: string(sen.Attr.DataType),
					States:   sen.Attr.States,
				},
			}
			j++
		}

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
		jsonDevices[i] = jsonDevice{
			ID:              device.ID,
			Address:         device.Address,
			AddressRequired: device.AddressRequired,
			Name:            device.Name,
			Description:     device.Description,
			ModelNumber:     device.ModelNumber,
			ModelName:       device.ModelName,
			SoftwareVersion: device.SoftwareVersion,
			ClientID:        "device_" + strconv.Itoa(i),
			Zones:           jsonZones,
			Sensors:         jsonSensors,
			ConnPool:        connPoolJSON,
			Auth:            authJSON,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(struct {
		Devices []jsonDevice `json:"devices"`
		Scenes  []jsonScene  `json:"scenes"`
	}{
		Devices: jsonDevices,
		Scenes:  jsonScenes,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiDiscoveryHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		discovererID := vars["discovererID"]

		discoverer := system.Extensions.FindDiscovererFromID(system, discovererID)
		if discoverer == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := discoverer.ScanDevices(system)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		writeDiscoveryResults(res, w)
	}
}

func apiFromStringDiscoveryHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		discovererID := vars["discovererID"]

		discoverer := system.Extensions.FindDiscovererFromID(system, discovererID)
		if discoverer == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(string(body))

		unquotedBody, err := strconv.Unquote(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		res, err := discoverer.FromString(unquotedBody)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("%+v\n", res)
		writeDiscoveryResults(res, w)
	}
}

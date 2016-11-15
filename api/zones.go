package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
	errExt "github.com/pkg/errors"
)

// RegisterZoneHandlers registers all of the zone specific API REST routes
func RegisterZoneHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/zones",
		apiListZonesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/zones",
		apiAddZoneHandler(s.systemSavePath, s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/zones/{id}/level",
		apiUpdateZoneLevelHandler(s.system)).Methods("PUT")
	r.HandleFunc("/api/v1/zones/{id}",
		apiUpdateZoneHandler(s.systemSavePath, s.system, s.recipeManager)).Methods("PUT")
}

func ZonesToJSON(zns map[string]*zone.Zone) []jsonZone {
	zones := make(zones, len(zns))
	var i int32
	for _, zone := range zns {
		zones[i] = jsonZone{
			Address:     zone.Address,
			ID:          zone.ID,
			Name:        zone.Name,
			Description: zone.Description,
			Type:        zone.Type.ToString(),
			Output:      zone.Output.ToString(),
			DeviceID:    zone.DeviceID,
		}
		i++
	}
	sort.Sort(zones)
	return zones
}

func apiUpdateZoneLevelHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//{cmd: "setLevel", value: 21, r: 0, g: 0, b: 0}
		var x struct {
			CMD   string  `json:"cmd"`
			Value float32 `json:"value"`
			R     byte    `json:"r"`
			G     byte    `json:"g"`
			B     byte    `json:"b"`
		}
		if err = json.Unmarshal(body, &x); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		zone, ok := system.Zones[vars["id"]]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid zone ID: %s", vars["id"]), w)
			return
		}

		switch x.CMD {
		case "setLevel":
			desc := fmt.Sprintf("Zone[%s] set level v:%f, r:%d, g:%d, b:%d", zone.Name, x.Value, x.R, x.G, x.B)
			err = system.Services.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.ZoneSetLevel{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
				Level: cmd.Level{
					Value: x.Value,
					R:     x.R,
					G:     x.G,
					B:     x.B,
				},
			}))
		case "turnOn":
			desc := fmt.Sprintf("Zone[%s] turn on", zone.Name)
			err = system.Services.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.ZoneTurnOn{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
			}))
		case "turnOff":
			desc := fmt.Sprintf("Zone[%s] turn off", zone.Name)
			err = system.Services.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.ZoneTurnOff{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
			}))
		default:
			respBadRequest(fmt.Sprintf("unknown zone command: %s", x.CMD), w)
			return
		}

		if err != nil {
			respErr(errExt.Wrap(err, "failed to set zone level"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiUpdateZoneHandler(
	savePath string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			respBadRequest(errExt.Wrap(err, "failed to read request body").Error(), w)
			return
		}

		var data jsonZone
		if err = json.Unmarshal(body, &data); err != nil {
			respBadRequest(errExt.Wrap(err, "failed to parse JSON in request body").Error(), w)
			return
		}

		zn, ok := system.Zones[data.ID]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid zone ID: %s", data.ID), w)
			return
		}

		updatedZn := &zone.Zone{
			ID:          data.ID,
			Address:     data.Address,
			Name:        data.Name,
			Description: data.Description,
			DeviceID:    data.DeviceID,
			Type:        zone.TypeFromString(data.Type),
			Output:      zone.OutputFromString(data.Output),
		}

		_, ok = system.Devices[data.DeviceID]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid device ID: %s", data.DeviceID), w)
			return
		}

		errors := updatedZn.Validate()
		if errors != nil {
			respValErr(&data, data.ID, errors, w)
			return
		}

		// Validated, set the fields
		zn.Name = updatedZn.Name
		zn.Description = updatedZn.Description
		zn.Address = updatedZn.Address
		zn.Type = updatedZn.Type
		zn.Output = updatedZn.Output

		//TODO: Support
		//zn.DeviceID = updatedZn.DeviceID

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			respErr(errExt.Wrap(err, "failed to save system to disk"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

func apiListZonesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(ZonesToJSON(system.Zones)); err != nil {
			respErr(errExt.Wrap(err, "failed to encode JSON"), w)
		}
	}
}

func apiAddZoneHandler(
	savePath string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			respBadRequest(errExt.Wrap(err, "failed to read request body").Error(), w)
			return
		}

		var data jsonZone
		if err = json.Unmarshal(body, &data); err != nil {
			respBadRequest(errExt.Wrap(err, "failed to unmarshal request JSON").Error(), w)
			return
		}

		z := &zone.Zone{
			ID:          data.ID,
			Address:     data.Address,
			Name:        data.Name,
			Description: data.Description,
			DeviceID:    data.DeviceID,
			Type:        zone.TypeFromString(data.Type),
			Output:      zone.OutputFromString(data.Output),
		}

		valErrs := z.Validate()
		if valErrs != nil {
			respValErr(&data, data.ID, valErrs, w)
			return
		}

		dev, ok := system.Devices[data.DeviceID]
		if !ok {
			respBadRequest(fmt.Sprintf("unknown device ID: %s", data.DeviceID), w)
			return
		}
		errors := dev.AddZone(z)
		if errors != nil {
			if valErrs, ok := errors.(*validation.Errors); ok {
				respValErr(&data, data.ID, valErrs, w)
			} else {
				respBadRequest(errExt.Wrap(err, "failed to add zone to device").Error(), w)
			}
			return
		}

		system.AddZone(z)

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			respErr(errExt.Wrap(err, "failed to save system to disk"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

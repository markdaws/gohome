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
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

// RegisterZoneHandlers registers all of the zone specific API REST routes
func RegisterZoneHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/zones",
		apiZonesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/zones",
		apiAddZoneHandler(s.systemSavePath, s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/zones/{id}",
		apiZoneHandler(s.system)).Methods("PUT")
}

func apiZoneHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch x.CMD {
		case "setLevel":
			desc := fmt.Sprintf("Zone[%s] set level v:%f, r:%d, g:%d, b:%d", zone.Name, x.Value, x.R, x.G, x.B)
			err = system.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.ZoneSetLevel{
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
			err = system.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.ZoneTurnOn{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
			}))
		case "turnOff":
			desc := fmt.Sprintf("Zone[%s] turn off", zone.Name)
			err = system.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.ZoneTurnOff{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
			}))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			log.E("failed to enqueue ZoneSetLevel command, ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiZonesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		zones := make(zones, len(system.Zones), len(system.Zones))
		var i int32
		for _, zone := range system.Zones {
			zones[i] = jsonZone{
				Address:     zone.Address,
				ID:          zone.ID,
				Name:        zone.Name,
				Description: zone.Description,
				Type:        zone.Type.ToString(),
				Output:      zone.Output.ToString(),
			}
			i++
		}
		sort.Sort(zones)
		if err := json.NewEncoder(w).Encode(zones); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data jsonZone
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		z := &zone.Zone{
			Address:     data.Address,
			Name:        data.Name,
			Description: data.Description,
			DeviceID:    data.DeviceID,
			Type:        zone.TypeFromString(data.Type),
			Output:      zone.OutputFromString(data.Output),
		}

		errors := system.AddZone(z)
		if errors != nil {
			if valErrs, ok := errors.(*validation.Errors); ok {
				fmt.Printf("%+v\n", valErrs.Errors[0])
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&data, data.ClientID, valErrs))
			} else {
				//Other kind of errors, TODO: log
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data.ClientID = ""
		data.ID = z.ID
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

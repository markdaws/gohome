package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
)

// RegisterSensorHandlers registers the REST API routes relating to devices
func RegisterSensorHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/sensors",
		apiSensorsHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/sensors",
		apiAddSensorHandler(s.systemSavePath, s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/sensors/{id}",
		apiUpdateSensorHandler(s.systemSavePath, s.system, s.recipeManager)).Methods("PUT")
}

func SensorsToJSON(sens map[string]*gohome.Sensor) []jsonSensor {
	sensors := make(sensors, len(sens))
	var i int32
	for _, sensor := range sens {
		sensors[i] = jsonSensor{
			Address:     sensor.Address,
			ID:          sensor.ID,
			Name:        sensor.Name,
			Description: sensor.Description,
			DeviceID:    sensor.DeviceID,
			Attr: jsonSensorAttr{
				Name:     sensor.Attr.Name,
				DataType: string(sensor.Attr.DataType),
				States:   sensor.Attr.States,
			},
		}
		i++
	}
	sort.Sort(sensors)
	return sensors
}

func apiSensorsHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(SensorsToJSON(system.Sensors)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiUpdateSensorHandler(
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

		var data jsonSensor
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sen, ok := system.Sensors[data.ID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sen.Name = data.Name
		sen.Address = data.Address
		sen.Description = data.Description

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

func apiAddSensorHandler(
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

		var data jsonSensor
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sensor := &gohome.Sensor{
			ID:          data.ID,
			Name:        data.Name,
			Description: data.Description,
			Address:     data.Address,
			DeviceID:    data.DeviceID,
			Attr: gohome.SensorAttr{
				Name:     data.Attr.Name,
				DataType: gohome.SensorDataType(data.Attr.DataType),
				States:   data.Attr.States,
			},
		}

		valErrs := sensor.Validate()
		if valErrs != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(validation.NewErrorJSON(&data, data.ID, valErrs))
			return
		}

		//Add the sensor to the owner device
		dev, ok := system.Devices[data.DeviceID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = dev.AddSensor(sensor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		system.AddSensor(sensor)

		err = store.SaveSystem(savePath, system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

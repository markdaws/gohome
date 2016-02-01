package www

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/discovery"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

type wwwServer struct {
	rootPath      string
	system        *gohome.System
	recipeManager *gohome.RecipeManager
	eventLogger   gohome.WSEventLogger
}

// ListenAndServe creates a new WWW server, that handles API calls and also
// runs the gohome website
func ListenAndServe(
	rootPath string,
	port string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager,
	eventLogger gohome.WSEventLogger) error {
	server := &wwwServer{
		rootPath:      rootPath,
		system:        system,
		recipeManager: recipeManager,
		eventLogger:   eventLogger,
	}
	return server.listenAndServe(port)
}

func (s *wwwServer) listenAndServe(port string) error {

	r := mux.NewRouter()

	mime.AddExtensionType(".jsx", "text/jsx")
	mime.AddExtensionType(".woff", "application/font-woff")
	mime.AddExtensionType(".woff2", "application/font-woff2")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")
	cssHandler := http.FileServer(http.Dir(s.rootPath + "/assets/css/"))
	jsHandler := http.FileServer(http.Dir(s.rootPath + "/assets/js/"))
	fontHandler := http.FileServer(http.Dir(s.rootPath + "/assets/fonts/"))
	jsxHandler := http.FileServer(http.Dir(s.rootPath + "/assets/jsx/"))
	imageHandler := http.FileServer(http.Dir(s.rootPath + "/assets/images/"))

	// Websocket handler
	r.HandleFunc("/api/v1/events/ws", s.eventLogger.HTTPHandler())

	r.HandleFunc("/api/v1/systems/{systemId}/scenes", apiScenesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/systems/{systemId}/scenes/{id}", apiSceneHandlerUpdate(s.system, s.recipeManager)).Methods("PUT")
	r.HandleFunc("/api/v1/systems/{systemId}/scenes/{sceneId}/commands/{index}", apiSceneHandlerCommandDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/systems/{systemId}/scenes/{sceneId}/commands", apiSceneHandlerCommandAdd(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/systems/{systemId}/scenes/{id}", apiSceneHandlerDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/systems/{systemId}/scenes/active", apiActiveScenesHandler(s.system)).Methods("POST")

	r.HandleFunc("/api/v1/systems/{systemId}/zones", apiZonesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/systems/{systemId}/zones", apiAddZoneHandler(s.system)).Methods("POST")
	r.HandleFunc("/api/v1/systems/{systemId}/zones/{id}", apiZoneHandler(s.system)).Methods("GET")

	r.HandleFunc("/api/v1/systems/{systemId}/devices", apiDevicesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/systems/{systemId}/devices", apiAddDeviceHandler(s.system)).Methods("POST")

	// Discover devices and capabilities on the network
	r.HandleFunc("/api/v1/discovery/{modelNumber}", apiDiscoveryHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/discovery/{modelNumber}/token", apiDiscoveryTokenHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/discovery/{modelNumber}/access", apiDiscoveryAccessHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/discovery/{modelNumber}/zones", apiDiscoveryZoneHandler(s.system)).Methods("GET")

	r.HandleFunc("/api/v1/cookbooks", apiCookBooksHandler(s.recipeManager.CookBooks)).Methods("GET")
	r.HandleFunc("/api/v1/cookbooks/{id}", apiCookBookHandler(s.recipeManager.CookBooks)).Methods("GET")

	//TODO: Need System in the api path?
	r.HandleFunc("/api/v1/recipes", apiRecipesHandlerPost(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/recipes/{id}", apiRecipeHandler(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/recipes/{id}", apiRecipeHandlerDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/recipes", apiRecipesHandlerGet(s.system, s.recipeManager)).Methods("GET")

	sub := r.PathPrefix("/assets").Subrouter()
	//sub.Methods("GET")
	sub.Handle("/css/{filename}", http.StripPrefix("/assets/css/", cssHandler))
	sub.Handle("/js/{filename}", http.StripPrefix("/assets/js/", jsHandler))
	sub.Handle("/fonts/{filename}", http.StripPrefix("/assets/fonts/", fontHandler))
	sub.Handle("/jsx/{filename}", http.StripPrefix("/assets/jsx/", jsxHandler))
	sub.Handle("/images/{filename}", http.StripPrefix("/assets/images/", imageHandler))
	r.HandleFunc("/", rootHandler(s.rootPath))
	return http.ListenAndServe(port, r)
}

func rootHandler(rootPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rootPath+"/assets/html/index.html")
	}
}

func apiRecipesHandlerPost(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var data map[string]interface{}
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipe, err := recipeManager.UnmarshalNewRecipe(data)
		if err != nil {
			errBad := err.(*gohome.ErrUnmarshalRecipe)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(struct {
				ParamID     string `json:"paramId"`
				ErrorType   string `json:"errorType"`
				Description string `json:"description"`
			}{
				ParamID:     errBad.ParamID,
				ErrorType:   errBad.ErrorType,
				Description: errBad.Description,
			})
			return
		}

		recipeManager.RegisterAndStart(recipe)
		err = system.Save(recipeManager)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{ID: recipe.ID})
	}
}

func apiRecipeHandler(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var data struct {
			Enabled bool `json:"enabled"`
		}
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeID := mux.Vars(r)["id"]
		recipe := recipeManager.RecipeByID(recipeID)
		if recipe == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = recipeManager.EnableRecipe(recipe, data.Enabled)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		system.Save(recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiRecipeHandlerDelete(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		recipeID := mux.Vars(r)["id"]
		recipe := recipeManager.RecipeByID(recipeID)
		if recipe == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := recipeManager.DeleteRecipe(recipe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		system.Save(recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiRecipesHandlerGet(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		recipes := system.Recipes
		jsonRecipes := make(jsonRecipes, len(recipes))

		i := 0
		for _, recipe := range recipes {
			jsonRecipes[i] = jsonRecipe{
				ID:          recipe.ID,
				Name:        recipe.Name,
				Description: recipe.Description,
				Enabled:     recipe.Enabled(),
			}
			i++
		}
		sort.Sort(jsonRecipes)
		if err := json.NewEncoder(w).Encode(jsonRecipes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiCookBooksHandler(cookBooks []*gohome.CookBook) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		type jsonCookBook struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			LogoURL     string `json:"logoUrl"`
		}

		//TODO: Return in a consistent order
		jsonCookBooks := make([]jsonCookBook, len(cookBooks))
		for i, cookBook := range cookBooks {
			jsonCookBooks[i] = jsonCookBook{
				ID:          cookBook.ID,
				Name:        cookBook.Name,
				Description: cookBook.Description,
				LogoURL:     cookBook.LogoURL,
			}
		}
		if err := json.NewEncoder(w).Encode(jsonCookBooks); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiCookBookHandler(cookBooks []*gohome.CookBook) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		type jsonIngredient struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Type        string `json:"type"`
			Reference   string `json:"reference"`
		}
		type jsonTrigger struct {
			ID          string           `json:"id"`
			Name        string           `json:"name"`
			Description string           `json:"description"`
			Ingredients []jsonIngredient `json:"ingredients"`
		}
		type jsonAction struct {
			ID          string           `json:"id"`
			Name        string           `json:"name"`
			Description string           `json:"description"`
			Ingredients []jsonIngredient `json:"ingredients"`
		}
		type jsonCookBook struct {
			ID          string        `json:"id"`
			Name        string        `json:"name"`
			Description string        `json:"description"`
			LogoURL     string        `json:"logoUrl"`
			Triggers    []jsonTrigger `json:"triggers"`
			Actions     []jsonAction  `json:"actions"`
		}

		vars := mux.Vars(r)
		cbID := vars["id"]
		var found = false
		for _, c := range cookBooks {
			if c.ID != cbID {
				continue
			}

			jsonTriggers := make([]jsonTrigger, len(c.Triggers))
			for i, t := range c.Triggers {
				jsonTriggers[i] = jsonTrigger{
					ID:          t.Type(),
					Name:        t.Name(),
					Description: t.Description(),
					Ingredients: make([]jsonIngredient, len(t.Ingredients())),
				}

				for j, ing := range t.Ingredients() {
					jsonTriggers[i].Ingredients[j] = jsonIngredient{
						ID:          ing.ID,
						Name:        ing.Name,
						Description: ing.Description,
						Type:        ing.Type,
					}
				}
			}

			// for each trigger need to json all ingredients
			jsonActions := make([]jsonAction, len(c.Actions))
			for i, a := range c.Actions {
				jsonActions[i] = jsonAction{
					ID:          a.Type(),
					Name:        a.Name(),
					Description: a.Description(),
					Ingredients: make([]jsonIngredient, len(a.Ingredients())),
				}

				for j, ing := range a.Ingredients() {
					jsonActions[i].Ingredients[j] = jsonIngredient{
						ID:          ing.ID,
						Name:        ing.Name,
						Description: ing.Description,
						Type:        ing.Type,
					}
				}
			}

			// for each action need to json all ingredients
			jsonCookBook := jsonCookBook{
				ID:          c.ID,
				Name:        c.Name,
				Description: c.Description,
				LogoURL:     c.LogoURL,
				Triggers:    jsonTriggers,
				Actions:     jsonActions,
			}
			if err := json.NewEncoder(w).Encode(jsonCookBook); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			found = true
			break
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func apiScenesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		scenes := make(scenes, len(system.Scenes), len(system.Scenes))
		var i int32
		for _, scene := range system.Scenes {
			scenes[i] = jsonScene{
				Address:     scene.Address,
				ID:          scene.ID,
				Name:        scene.Name,
				Description: scene.Description,
				Managed:     scene.Managed,
			}

			cmds := make([]jsonCommand, len(scene.Commands))
			for j, sCmd := range scene.Commands {
				switch xCmd := sCmd.(type) {
				case *cmd.ZoneSetLevel:
					cmds[j] = jsonCommand{
						Type: "zoneSetLevel",
						Attributes: map[string]interface{}{
							"ZoneID": xCmd.ZoneID,
							"Level":  xCmd.Level.Value,
						},
					}
				case *cmd.ButtonPress:
					cmds[j] = jsonCommand{
						Type: "buttonPress",
						Attributes: map[string]interface{}{
							"ButtonID": xCmd.ButtonID,
						},
					}
				case *cmd.ButtonRelease:
					cmds[j] = jsonCommand{
						Type: "buttonRelease",
						Attributes: map[string]interface{}{
							"ButtonID": xCmd.ButtonID,
						},
					}
				case *cmd.SceneSet:
					cmds[j] = jsonCommand{
						Type: "sceneSet",
						Attributes: map[string]interface{}{
							"SceneID": xCmd.SceneID,
						},
					}
				default:
					fmt.Println("unknown scene command")
				}
			}

			scenes[i].Commands = cmds
			i++
		}
		sort.Sort(scenes)
		if err := json.NewEncoder(w).Encode(scenes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiSceneHandlerDelete(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sceneID := mux.Vars(r)["id"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		system.DeleteScene(scene)
		err := system.Save(recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerCommandDelete(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sceneID := mux.Vars(r)["sceneId"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		commandIndex, err := strconv.Atoi(mux.Vars(r)["index"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = scene.DeleteCommand(commandIndex)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = system.Save(recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerCommandAdd(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		sceneID := mux.Vars(r)["sceneId"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var command jsonCommand
		if err = json.Unmarshal(body, &command); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var finalCmd cmd.Command
		switch command.Type {
		case "zoneSetLevel":
			if _, ok := command.Attributes["ZoneID"]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attribute_ZoneID", "required field", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			if _, ok = command.Attributes["ZoneID"].(string); !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attributes_ZoneID", "must be a string data type", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			z, ok := system.Zones[command.Attributes["ZoneID"].(string)]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				var valErrs *validation.Errors
				if command.Attributes["ZoneID"].(string) == "" {
					valErrs = validation.NewErrors("attributes_ZoneID", "required field", true)
				} else {
					valErrs = validation.NewErrors("attributes_ZoneID", "invalid zone ID", true)
				}
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			_, ok = command.Attributes["Level"]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attribute_Level", "required field", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}
			if _, ok = command.Attributes["Level"].(float64); !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attribute_Level", "must be a float data type", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			finalCmd = &cmd.ZoneSetLevel{
				ZoneAddress: z.Address,
				ZoneID:      z.ID,
				ZoneName:    z.Name,
				Level:       cmd.Level{Value: float32(command.Attributes["Level"].(float64))},
			}
		case "buttonPress":
			//TODO:
			w.WriteHeader(http.StatusBadRequest)
			return
		case "buttonRelease":
			//TODO:
			w.WriteHeader(http.StatusBadRequest)
			return
		case "sceneSet":
			if _, ok := command.Attributes["SceneID"]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attribute_SceneID", "required field", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			if _, ok = command.Attributes["SceneID"].(string); !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attributes_SceneID", "must be a string data type", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			scene, ok := system.Scenes[command.Attributes["SceneID"].(string)]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				var valErrs *validation.Errors
				if command.Attributes["SceneID"].(string) == "" {
					valErrs = validation.NewErrors("attributes_SceneID", "required field", true)
				} else {
					valErrs = validation.NewErrors("attributes_SceneID", "invalid Scene ID", true)
				}
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}
			finalCmd = &cmd.SceneSet{
				scene.ID,
				scene.Name,
			}

		default:
			//TODO:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = scene.AddCommand(finalCmd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = system.Save(recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerUpdate(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//TODO
		/*
			sceneID := mux.Vars(r)["id"]
			scene, ok := system.Scenes[sceneID]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// What actions can the user perform
			// Add command
			// Delete command
			// general name update
			// command update


			system.DeleteScene(scene)
			err := system.Save(recipeManager)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(struct{}{})
		*/
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
				Controller:  zone.Controller,
			}
			i++
		}
		sort.Sort(zones)
		if err := json.NewEncoder(w).Encode(zones); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiAddZoneHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
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
			Controller:  data.Controller,
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

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiDevicesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		devices := make(devices, len(system.Devices), len(system.Devices))
		var i int32
		for _, device := range system.Devices {
			devices[i] = jsonDevice{
				Address:     device.Address(),
				ID:          device.ID(),
				Name:        device.Name(),
				Description: device.Description(),
				ModelNumber: device.ModelNumber(),
			}
			i++
		}
		sort.Sort(devices)
		if err := json.NewEncoder(w).Encode(devices); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiAddDeviceHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
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

		var auth *comm.Auth
		if data.Token != "" {
			auth = &comm.Auth{
				Token: data.Token,
			}
		}
		//TODO: Don't pass in ID
		d := gohome.NewDevice(
			data.ModelNumber,
			data.Address,
			system.NextGlobalID(),
			data.Name,
			data.Description,
			false, //TODO: stream?
			auth,
		)

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

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
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
			err = system.CmdProcessor.Enqueue(&cmd.ZoneSetLevel{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
				Level: cmd.Level{
					Value: x.Value,
					R:     x.R,
					G:     x.G,
					B:     x.B,
				},
			})
		case "turnOn":
			err = system.CmdProcessor.Enqueue(&cmd.ZoneTurnOn{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
			})
		case "turnOff":
			err = system.CmdProcessor.Enqueue(&cmd.ZoneTurnOff{
				ZoneAddress: zone.Address,
				ZoneID:      zone.ID,
				ZoneName:    zone.Name,
			})
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiActiveScenesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var x struct {
			ID string `json:"id"`
		}
		if err = json.Unmarshal(body, &x); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		scene, ok := system.Scenes[x.ID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = system.CmdProcessor.Enqueue(&cmd.SceneSet{
			SceneID:   scene.ID,
			SceneName: scene.Name,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiDiscoveryHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//This is blocking
		data, err := discovery.Discover(vars["modelNumber"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Location string `json:"location"`
		}{Location: data["location"]})
	}
}

func apiDiscoveryTokenHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//This is blocking
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

		//This is blocking
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
				Controller:  zone.Controller,
			}
			i++
		}
		sort.Sort(jsonZones)
		if err := json.NewEncoder(w).Encode(jsonZones); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

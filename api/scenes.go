package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

// RegisterSceneHandlers registers all of the scene specific API REST routes
func RegisterSceneHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/scenes",
		apiScenesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/scenes/{id}",
		apiSceneHandlerUpdate(s.system, s.recipeManager)).Methods("PUT")
	r.HandleFunc("/api/v1/scenes",
		apiSceneHandlerCreate(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/scenes/{sceneId}/commands/{index}",
		apiSceneHandlerCommandDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/scenes/{sceneId}/commands",
		apiSceneHandlerCommandAdd(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/scenes/{id}",
		apiSceneHandlerDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/scenes/active",
		apiActiveScenesHandler(s.system)).Methods("POST")
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

		desc := fmt.Sprintf("Set scene: %s", scene.Name)
		err = system.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.SceneSet{
			SceneID:   scene.ID,
			SceneName: scene.Name,
		}))
		if err != nil {
			//TODO: log
			fmt.Printf("enqueue failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
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
		err := store.SaveSystem(system, recipeManager)
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

		err = store.SaveSystem(system, recipeManager)
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
				valErrs := validation.NewErrors("attributes_Level", "required field", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}
			if _, ok = command.Attributes["Level"].(float64); !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attributes_Level", "must be a float data type", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			var r, g, b byte
			if z.Output == zone.OTRGB {
				_, ok = command.Attributes["R"]
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					valErrs := validation.NewErrors("attributes_R", "required field", true)
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
					return
				}
				if _, ok = command.Attributes["R"].(float64); !ok {
					w.WriteHeader(http.StatusBadRequest)
					valErrs := validation.NewErrors("attributes_R", "must be an integer data type", true)
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
					return
				}

				_, ok = command.Attributes["G"]
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					valErrs := validation.NewErrors("attributes_G", "required field", true)
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
					return
				}
				if _, ok = command.Attributes["G"].(float64); !ok {
					w.WriteHeader(http.StatusBadRequest)
					valErrs := validation.NewErrors("attributes_G", "must be an integer data type", true)
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
					return
				}

				_, ok = command.Attributes["B"]
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					valErrs := validation.NewErrors("attributes_B", "required field", true)
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
					return
				}
				if _, ok = command.Attributes["B"].(float64); !ok {
					w.WriteHeader(http.StatusBadRequest)
					valErrs := validation.NewErrors("attributes_B", "must be an integer data type", true)
					json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
					return
				}

				r = byte(command.Attributes["R"].(float64))
				g = byte(command.Attributes["G"].(float64))
				b = byte(command.Attributes["B"].(float64))
			}

			finalCmd = &cmd.ZoneSetLevel{
				ZoneAddress: z.Address,
				ZoneID:      z.ID,
				ZoneName:    z.Name,
				Level: cmd.Level{
					Value: float32(command.Attributes["Level"].(float64)),
					R:     r,
					G:     g,
					B:     b,
				},
			}
		case "buttonPress", "buttonRelease":
			if _, ok := command.Attributes["ButtonID"]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attribute_ButtonID", "required field", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			if _, ok = command.Attributes["ButtonID"].(string); !ok {
				w.WriteHeader(http.StatusBadRequest)
				valErrs := validation.NewErrors("attributes_ButtonID", "must be a string data type", true)
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			button, ok := system.Buttons[command.Attributes["ButtonID"].(string)]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				var valErrs *validation.Errors
				if command.Attributes["ButtonID"].(string) == "" {
					valErrs = validation.NewErrors("attributes_ButtonID", "required field", true)
				} else {
					valErrs = validation.NewErrors("attributes_ButtonID", "invalid Button ID", true)
				}
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&command, command.ClientID, valErrs))
				return
			}

			if command.Type == "buttonPress" {
				finalCmd = &cmd.ButtonPress{
					ButtonAddress: button.Address,
					ButtonID:      button.ID,
					DeviceName:    button.Device.Name,
					DeviceAddress: button.Device.Address,
					DeviceID:      button.Device.ID,
				}
			} else {
				finalCmd = &cmd.ButtonRelease{
					ButtonAddress: button.Address,
					ButtonID:      button.ID,
					DeviceName:    button.Device.Name,
					DeviceAddress: button.Device.Address,
					DeviceID:      button.Device.ID,
				}
			}

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = scene.AddCommand(finalCmd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = store.SaveSystem(system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerUpdate(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sceneID := mux.Vars(r)["id"]
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

		var updates struct {
			Name        *string `json:"name"`
			Address     *string `json:"address"`
			Description *string `json:"description"`
		}
		if err = json.Unmarshal(body, &updates); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var updatedScene = *scene
		if updates.Name != nil {
			updatedScene.Name = *updates.Name
		}
		if updates.Address != nil {
			updatedScene.Address = *updates.Address
		}
		if updates.Description != nil {
			updatedScene.Description = *updates.Description
		}

		valErrs := updatedScene.Validate()
		if valErrs != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(validation.NewErrorJSON(&updates, sceneID, valErrs))
			return
		}

		system.Scenes[sceneID] = &updatedScene

		err = store.SaveSystem(system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerCreate(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var scene jsonScene
		if err = json.Unmarshal(body, &scene); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newScene := &gohome.Scene{
			Address:     scene.Address,
			Name:        scene.Name,
			Description: scene.Description,
			Managed:     true,
		}

		err = system.AddScene(newScene)
		if err != nil {
			if valErrs, ok := err.(*validation.Errors); ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(validation.NewErrorJSON(&scene, scene.ClientID, valErrs))
			} else {
				//Other kind of errors, TODO: log
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		err = store.SaveSystem(system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		scene.ID = newScene.ID
		scene.ClientID = ""
		json.NewEncoder(w).Encode(scene)
	}
}

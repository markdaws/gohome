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
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/validation"
	errExt "github.com/pkg/errors"
)

// RegisterSceneHandlers registers all of the scene specific API REST routes
func RegisterSceneHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/scenes", apiScenesHandler(s.system)).Methods("GET")

	r.HandleFunc("/api/v1/scenes/{ID}",
		apiSceneHandlerUpdate(s.systemSavePath, s.system)).Methods("PUT")

	r.HandleFunc("/api/v1/scenes",
		apiSceneHandlerCreate(s.systemSavePath, s.system)).Methods("POST")

	r.HandleFunc("/api/v1/scenes/{sceneID}/commands/{commandID}",
		apiSceneHandlerCommandDelete(s.systemSavePath, s.system)).Methods("DELETE")

	r.HandleFunc("/api/v1/scenes/{sceneID}/commands",
		apiSceneHandlerCommandAdd(s.systemSavePath, s.system)).Methods("POST")

	r.HandleFunc("/api/v1/scenes/{ID}",
		apiSceneHandlerDelete(s.systemSavePath, s.system)).Methods("DELETE")

	r.HandleFunc("/api/v1/scenes/active",
		apiActiveScenesHandler(s.system)).Methods("POST")
}

func ScenesToJSON(inputScenes map[string]*gohome.Scene) scenes {
	jsonScenes := make(scenes, len(inputScenes))
	var i int32
	for _, scene := range inputScenes {
		jsonScenes[i] = jsonScene{
			Address:     scene.Address,
			ID:          scene.ID,
			Name:        scene.Name,
			Description: scene.Description,
			Managed:     scene.Managed,
		}

		cmds := make([]jsonCommand, len(scene.Commands))
		for j, sCmd := range scene.Commands {
			switch xCmd := sCmd.(type) {
			case *cmd.SceneSet:
				cmds[j] = jsonCommand{
					ID:   xCmd.ID,
					Type: "sceneSet",
					Attributes: map[string]interface{}{
						"id":      xCmd.ID,
						"SceneID": xCmd.SceneID,
					},
				}

			case *cmd.FeatureSetAttrs:
				cmds[j] = jsonCommand{
					ID:   xCmd.ID,
					Type: "featureSetAttrs",
					Attributes: map[string]interface{}{
						"id":    xCmd.FeatureID,
						"type":  xCmd.FeatureType,
						"attrs": xCmd.Attrs,
					},
				}
			default:
				fmt.Println("unknown scene command")
			}
		}

		jsonScenes[i].Commands = cmds
		i++
	}
	sort.Sort(jsonScenes)
	return jsonScenes
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
		err = system.Services.CmdProcessor.Enqueue(gohome.NewCommandGroup(desc, &cmd.SceneSet{
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

		if err := json.NewEncoder(w).Encode(ScenesToJSON(system.Scenes)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiSceneHandlerDelete(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sceneID := mux.Vars(r)["ID"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		system.DeleteScene(scene)
		err := store.SaveSystem(savePath, system)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerCommandDelete(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sceneID := mux.Vars(r)["sceneID"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid scene ID: %s", sceneID), w)
			return
		}

		commandID := mux.Vars(r)["commandID"]
		var command cmd.Command
		for _, c := range scene.Commands {
			if c.GetID() == commandID {
				command = c
				break
			}
		}
		if command == nil {
			respBadRequest(fmt.Sprintf("invalid command ID: %s", commandID), w)
			return
		}

		err := scene.DeleteCommand(commandID)
		if err != nil {
			respBadRequest(fmt.Sprintf("scene: %s does not contain command: %s", sceneID, commandID), w)
			return
		}

		err = store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "error writing changes to disk"), w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiSceneHandlerCommandAdd(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		sceneID := mux.Vars(r)["sceneID"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid scene ID: %s", sceneID), w)
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
		if err != nil {
			respBadRequest("unable to read request body", w)
			return
		}

		var command map[string]*json.RawMessage
		if err = json.Unmarshal(body, &command); err != nil {
			respBadRequest(errExt.Wrap(err, "unable to parse JSON in request body").Error(), w)
			return
		}

		var cmdType string
		if err = json.Unmarshal(*command["type"], &cmdType); err != nil {
			respBadRequest("invalid JSON body, missing type field", w)
			return
		}

		var finalCmd cmd.Command
		switch cmdType {
		case "featureSetAttrs":
			var cmdAttrs map[string]*json.RawMessage
			if err = json.Unmarshal(*command["attributes"], &cmdAttrs); err != nil {
				respBadRequest("invalid JSON body, missing attributes key", w)
				return
			}

			var featureID string
			if err = json.Unmarshal(*cmdAttrs["id"], &featureID); err != nil {
				respBadRequest("invalid JSON body, missing ID key", w)
				return
			}

			f, ok := system.Features[featureID]
			if !ok {
				respBadRequest(fmt.Sprintf("invalid feature ID: %s", featureID), w)
				return
			}

			attrs := make(map[string]*attr.Attribute)
			if err = json.Unmarshal(*cmdAttrs["attrs"], &attrs); err != nil {
				respBadRequest("invalid JSON body, missing attrs key", w)
				return
			}
			attr.FixJSON(attrs)

			finalCmd = &cmd.FeatureSetAttrs{
				ID:          system.NewGlobalID(),
				FeatureID:   featureID,
				FeatureName: f.Name,
				FeatureType: f.Type,
				Attrs:       attrs,
			}

		case "sceneSet":
			var sceneCmd jsonCommand
			if err = json.Unmarshal(body, &sceneCmd); err != nil {
				respBadRequest("invalid JSON in request body", w)
				return
			}

			if _, ok := sceneCmd.Attributes["SceneID"]; !ok {
				valErrs := validation.NewErrors("attribute_SceneID", "required field", true)
				respValErr(&sceneCmd, sceneCmd.ID, valErrs, w)
				return
			}

			scene, ok := system.Scenes[sceneCmd.Attributes["SceneID"].(string)]
			if !ok {
				var valErrs *validation.Errors
				if sceneCmd.Attributes["SceneID"].(string) == "" {
					valErrs = validation.NewErrors("attributes_SceneID", "required field", true)
				} else {
					valErrs = validation.NewErrors("attributes_SceneID", "invalid Scene ID", true)
				}
				respValErr(&sceneCmd, sceneCmd.ID, valErrs, w)
				return
			}
			finalCmd = &cmd.SceneSet{
				ID:        system.NewGlobalID(),
				SceneID:   scene.ID,
				SceneName: scene.Name,
			}

		default:
			respBadRequest(fmt.Sprintf("invalid command in type field: %s", cmdType), w)
			return
		}

		cmdID := finalCmd.GetID()
		err = scene.AddCommand(finalCmd)
		if err != nil {
			respBadRequest(errExt.Wrap(err, "failed to add command to scene").Error(), w)
			return
		}

		err = store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "error writing changes to disk"), w)
			return
		}

		// Need to send back the command with its new ID
		var updatedCmd jsonCommand
		if err = json.Unmarshal(body, &updatedCmd); err != nil {
			respBadRequest(errExt.Wrap(err, "unable to parse JSON in request body").Error(), w)
			return
		}
		updatedCmd.ID = cmdID
		json.NewEncoder(w).Encode(updatedCmd)
	}
}

func apiSceneHandlerUpdate(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sceneID := mux.Vars(r)["ID"]
		scene, ok := system.Scenes[sceneID]
		if !ok {
			respBadRequest(fmt.Sprintf("invalid scene ID: %s", sceneID), w)
			return
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			respBadRequest("unable to read request body", w)
			return
		}

		var updates struct {
			Name        *string `json:"name"`
			Address     *string `json:"address"`
			Description *string `json:"description"`
		}
		if err = json.Unmarshal(body, &updates); err != nil {
			respBadRequest("unable to parse request body, invalid JSON", w)
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
			respValErr(&updates, sceneID, valErrs, w)
			return
		}

		system.Scenes[sceneID] = &updatedScene

		err = store.SaveSystem(savePath, system)
		if err != nil {
			respErr(errExt.Wrap(err, "failed to save system to disk"), w)
			return
		}

		// TODO: Jsonify scene object and return full object, not fields user passed in
		// same for devices
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(updates)
	}
}

func apiSceneHandlerCreate(savePath string, system *gohome.System) func(http.ResponseWriter, *http.Request) {
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
			ID:          system.NewGlobalID(),
			Address:     scene.Address,
			Name:        scene.Name,
			Description: scene.Description,
			Managed:     true,
		}
		valErrs := newScene.Validate()
		if valErrs != nil {
			respValErr(&scene, scene.ID, valErrs, w)
			return
		}
		system.AddScene(newScene)

		err = store.SaveSystem(savePath, system)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		scene.ID = newScene.ID
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(scene)
	}
}

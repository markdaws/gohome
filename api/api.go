package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

type apiServer struct {
	system        *gohome.System
	recipeManager *gohome.RecipeManager
	eventLogger   gohome.WSEventLogger
}

// ListenAndServe creates a new WWW server, that handles API calls and also
// runs the gohome website
func ListenAndServe(
	port string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager,
	eventLogger gohome.WSEventLogger) error {
	server := &apiServer{
		system:        system,
		recipeManager: recipeManager,
		eventLogger:   eventLogger,
	}
	return server.listenAndServe(port)
}

//TODO: Clean all these up and unify naming

func (s *apiServer) listenAndServe(port string) error {

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/events/ws", s.eventLogger.HTTPHandler())

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

	r.HandleFunc("/api/v1/buttons",
		apiButtonsHandler(s.system)).Methods("GET")

	r.HandleFunc("/api/v1/zones",
		apiZonesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/zones",
		apiAddZoneHandler(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/zones/{id}",
		apiZoneHandler(s.system)).Methods("PUT")

	r.HandleFunc("/api/v1/devices",
		apiDevicesHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/devices",
		apiAddDeviceHandler(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/devices/{id}",
		apiDeviceHandlerDelete(s.system, s.recipeManager)).Methods("DELETE")

	r.HandleFunc("/api/v1/discovery/{modelNumber}",
		apiDiscoveryHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/discovery/{modelNumber}/token",
		apiDiscoveryTokenHandler(s.system)).Methods("GET")
	r.HandleFunc("/api/v1/discovery/{modelNumber}/access",
		apiDiscoveryAccessHandler(s.system)).Methods("GET")

	r.HandleFunc("/api/v1/cookbooks",
		apiCookBooksHandler(s.recipeManager.CookBooks)).Methods("GET")
	r.HandleFunc("/api/v1/cookbooks/{id}",
		apiCookBookHandler(s.recipeManager.CookBooks)).Methods("GET")

	r.HandleFunc("/api/v1/recipes",
		apiRecipesHandlerPost(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/recipes/{id}",
		apiRecipeHandler(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/recipes/{id}",
		apiRecipeHandlerDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/recipes",
		apiRecipesHandlerGet(s.system, s.recipeManager)).Methods("GET")

	return http.ListenAndServe(
		port,
		handlers.CORS(
			handlers.AllowedMethods([]string{"PUT", "POST", "DELETE", "GET", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)(r))
}

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/validation"
)

type apiServer struct {
	systemSavePath string
	system         *gohome.System
	recipeManager  *gohome.RecipeManager
	eventLogger    gohome.WSEventLogger
}

// ListenAndServe creates a new WWW server, that handles API calls and also
// runs the gohome website
func ListenAndServe(
	systemSavePath string,
	addr string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager,
	eventLogger gohome.WSEventLogger) error {
	server := &apiServer{
		systemSavePath: systemSavePath,
		system:         system,
		recipeManager:  recipeManager,
		eventLogger:    eventLogger,
	}
	return server.listenAndServe(addr)
}

func (s *apiServer) listenAndServe(addr string) error {

	r := mux.NewRouter()

	//TODO: re-enable
	//r.HandleFunc("/api/v1/events/ws", s.eventLogger.HTTPHandler())

	RegisterSceneHandlers(r, s)
	RegisterButtonHandlers(r, s)
	RegisterZoneHandlers(r, s)
	RegisterDeviceHandlers(r, s)
	RegisterSensorHandlers(r, s)
	RegisterDiscoveryHandlers(r, s)
	RegisterCookBookHandlers(r, s)
	RegisterRecipeHandlers(r, s)
	RegisterMonitorHandlers(r, s)

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: handlers.CORS(
			handlers.AllowedMethods([]string{"PUT", "POST", "DELETE", "GET", "OPTIONS", "UPGRADE"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)(r),
	}
	return server.ListenAndServe()
}

// apiResponse encapsulates the response from a http handler, responses can either
// be an error, such as invalid input, or contains a sucessful response
type apiResponse struct {
	// Err - will be non nil if there was a error processing the API request
	Err error

	// Data - pointer to struct that can be serialized to JSON that will then
	// be sent back to the client
	Data interface{}
}

// badRequestErr - API input was incorrect, e.g. missing required field.  The Msg field
// contains more specific information about the error
type badRequestErr struct {
	Msg string
}

func (r *badRequestErr) Error() string {
	return r.Msg
}

// validationErr - an error that occurs when input fields are not valid e.g. Name field
// is too long etc.
type validationErr struct {
	ID     string
	Data   interface{}
	Errors *validation.Errors
}

func (e *validationErr) Error() string {
	return e.Errors.Error()
}

// respBadRequest responds to the client with a http.StatusBadRequest and additional message
func respBadRequest(msg string, w http.ResponseWriter) {
	resp(apiResponse{
		Err: &badRequestErr{
			Msg: msg,
		},
	}, w)
}

func resp(r apiResponse, w http.ResponseWriter) {
	if r.Err != nil {
		switch err := r.Err.(type) {
		case *validationErr:
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(validation.NewErrorJSON(err.Data, err.ID, err.Errors))
		case *badRequestErr:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				Msg string
			}{Msg: err.Msg})
		default:
			log.E("Unknown error", r.Err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if r.Data != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(r.Data)
	}
}

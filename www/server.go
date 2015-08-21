package www

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

type Server interface {
	ListenAndServe(port string) error
}

type wwwServer struct {
	rootPath string
	system   *gohome.System
}

func NewServer(rootPath string, system *gohome.System) Server {
	return &wwwServer{rootPath: rootPath, system: system}
}

func (s *wwwServer) ListenAndServe(port string) error {
	r := mux.NewRouter()

	mime.AddExtensionType(".jsx", "text/jsx")
	cssHandler := http.FileServer(http.Dir(s.rootPath + "/assets/css/"))
	jsHandler := http.FileServer(http.Dir(s.rootPath + "/assets/js/"))
	jsxHandler := http.FileServer(http.Dir(s.rootPath + "/assets/jsx/"))
	imageHandler := http.FileServer(http.Dir(s.rootPath + "/assets/images/"))

	//TODO: Move api into separate http server
	r.HandleFunc("/api/v1/systems/{systemId}/scenes", apiScenesHandler(s.system))

	//TODO: Make for POST only
	//TODO: Have GET version to see the currently active scenes
	r.HandleFunc("/api/v1/systems/{systemId}/scenes/active", apiActiveScenesHandler(s.system)).Methods("POST")

	sub := r.PathPrefix("/assets").Subrouter()
	//sub.Methods("GET")
	sub.Handle("/css/{filename}", http.StripPrefix("/assets/css/", cssHandler))
	sub.Handle("/js/{filename}", http.StripPrefix("/assets/js/", jsHandler))
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

func apiScenesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		//TODO: Explicitly pull out required fields
		if err := json.NewEncoder(w).Encode(system.Scenes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiActiveScenesHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := r.Body.Close(); err != nil {
			//TODO: Error
			fmt.Println("Can't close body")
		}

		var x struct {
			Id string `json:"id"`
		}
		if err = json.Unmarshal(body, &x); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var foundScene bool = false
		for _, scene := range system.Scenes {
			if scene.Id == x.Id {
				foundScene = true
				scene.Execute()
				break
			}
		}
		if !foundScene {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

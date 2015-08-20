package www

import (
	"encoding/json"
	"fmt"
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
	r.HandleFunc("/api/systems/{systemId}/scenes", apiScenesHandler(s.system))

	r.Handle("/assets/css/{filename}", http.StripPrefix("/assets/css/", cssHandler))
	r.Handle("/assets/js/{filename}", http.StripPrefix("/assets/js/", jsHandler))
	r.Handle("/assets/jsx/{filename}", http.StripPrefix("/assets/jsx/", jsxHandler))
	r.Handle("/assets/images/{filename}", http.StripPrefix("/assets/images/", imageHandler))
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
		//loop through all scenes in the system, turn to JSON ...

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		if err := json.NewEncoder(w).Encode(system.Scenes); err != nil {
			fmt.Println("Failed to encode")
		}
	}
}

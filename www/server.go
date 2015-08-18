package www

import (
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

type Server interface {
	ListenAndServe(port string) error
}

type wwwServer struct {
	rootPath string
}

func NewServer(rootPath string) Server {
	return &wwwServer{rootPath: rootPath}
}

func (s *wwwServer) ListenAndServe(port string) error {
	r := mux.NewRouter()

	mime.AddExtensionType(".jsx", "text/jsx")
	cssHandler := http.FileServer(http.Dir(s.rootPath + "/assets/css/"))
	jsHandler := http.FileServer(http.Dir(s.rootPath + "/assets/js/"))
	jsxHandler := http.FileServer(http.Dir(s.rootPath + "/assets/jsx/"))
	imageHandler := http.FileServer(http.Dir(s.rootPath + "/assets/images/"))

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

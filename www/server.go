package www

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	gzip "github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

type wwwServer struct {
	rootPath string
	system   *gohome.System
	sessions *gohome.Sessions
}

// ListenAndServe creates a new WWW server, that handles API calls and also
// runs the gohome website
func ListenAndServe(
	rootPath string,
	addr string,
	system *gohome.System,
	sessions *gohome.Sessions) error {
	server := &wwwServer{
		rootPath: rootPath,
		system:   system,
		sessions: sessions,
	}
	return server.listenAndServe(addr)
}

func (s *wwwServer) listenAndServe(addr string) error {

	r := mux.NewRouter()

	mime.AddExtensionType(".jsx", "text/jsx")
	mime.AddExtensionType(".woff", "application/font-woff")
	mime.AddExtensionType(".woff2", "application/font-woff2")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")

	cssHandler := http.FileServer(http.Dir(s.rootPath + "/assets/css/"))
	extCssHandler := http.FileServer(http.Dir(s.rootPath + "/assets/css/ext/"))
	jsHandler := http.FileServer(http.Dir(s.rootPath + "/assets/js/"))
	jsExtHandler := http.FileServer(http.Dir(s.rootPath + "/assets/js/ext/"))
	fontHandler := http.FileServer(http.Dir(s.rootPath + "/assets/fonts/"))
	jsxHandler := http.FileServer(http.Dir(s.rootPath + "/assets/jsx/"))
	extImageHandler := http.FileServer(http.Dir(s.rootPath + "/assets/images/ext/"))
	imageHandler := http.FileServer(http.Dir(s.rootPath + "/assets/images/"))

	sub := r.PathPrefix("/assets").Subrouter()
	sub.Handle("/css/{filename}", http.StripPrefix("/assets/css/", gzip.GzipHandler(cssHandler)))
	sub.Handle("/css/ext/{filename}", http.StripPrefix("/assets/css/ext/", gzip.GzipHandler(extCssHandler)))
	sub.Handle("/js/{filename}", http.StripPrefix("/assets/js/", gzip.GzipHandler(jsHandler)))
	sub.Handle("/js/ext/{filename}", http.StripPrefix("/assets/js/ext/", gzip.GzipHandler(jsExtHandler)))
	sub.Handle("/fonts/{filename}", http.StripPrefix("/assets/fonts/", fontHandler))
	sub.Handle("/jsx/{filename}", http.StripPrefix("/assets/jsx/", gzip.GzipHandler(jsxHandler)))
	sub.Handle("/images/ext/{filename}", http.StripPrefix("/assets/images/ext/", extImageHandler))
	sub.Handle("/images/{filename}", http.StripPrefix("/assets/images/", imageHandler))

	r.HandleFunc("/api/v1/users/{login}/sessions", apiNewSessionHandler(s.system, s.sessions)).Methods("POST")
	r.HandleFunc("/logout", logoutHandler(s.rootPath))
	r.HandleFunc("/", rootHandler(s.rootPath))

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	return server.ListenAndServe()
}

func rootHandler(rootPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rootPath+"/assets/html/index.html")
	}
}

func logoutHandler(rootPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rootPath+"/assets/html/logout.html")
	}
}

func apiNewSessionHandler(sys *gohome.System, sessions *gohome.Sessions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var success = false
		login := mux.Vars(r)["login"]

		defer func() {
			sys.Services.EvtBus.Enqueue(&gohome.UserLoginEvt{
				Login:   login,
				Success: success,
			})
		}()

		var user *gohome.User
		for _, u := range sys.Users {
			if strings.ToLower(u.Login) == strings.ToLower(login) {
				user = u
				break
			}
		}
		if user == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var x struct {
			Password string `json:"password"`
		}
		if err = json.Unmarshal(body, &x); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = user.VerifyPassword(x.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sid, err := sessions.Add()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = sessions.Save()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:    "sid",
			Value:   sid,
			Path:    "/",
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)

		success = true

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			SessionID string `json:"sid"`
		}{
			SessionID: sid,
		})
	}
}

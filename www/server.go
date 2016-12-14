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

func cacheHandler(prefix string, h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if prefix != "" {
			temp := strings.Replace(r.URL.Path, prefix, "", -1)
			index := strings.Index(temp, "/")
			URLNoTimestamp := prefix + temp[index+1:]
			r.URL.Path = URLNoTimestamp
		}

		// Since all resources include cache busting values in the URL each time we have a new build,
		// we just set the max cache values here
		w.Header().Add("Cache-Control", "public; max-age=31536000")
		w.Header().Add("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
		h.ServeHTTP(w, r)
	}
}

func (s *wwwServer) listenAndServe(addr string) error {

	r := mux.NewRouter()

	mime.AddExtensionType(".jsx", "text/jsx")
	mime.AddExtensionType(".woff", "application/font-woff")
	mime.AddExtensionType(".woff2", "application/font-woff2")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")

	fileHandler := http.FileServer(http.Dir(s.rootPath + "/dist/"))

	sub := r.PathPrefix("/").Subrouter()

	// For each type of asset, they can be accessed via a direct url like /images/foo.png or they can
	// be accessed with some cache busting value prepended before the filename e.g.
	// /images/12345/foo.png which will redirect to foo.png on the filesystem. This allows you to either
	// put a hash value in the file name of an asset or some cache busting value like the build time in the
	// URL instead of having to rename files
	sub.HandleFunc("/js/{filename}", cacheHandler("", gzip.GzipHandler(fileHandler)))
	sub.HandleFunc("/js/{timestamp}/{filename}", cacheHandler("/js/", gzip.GzipHandler(fileHandler)))
	sub.HandleFunc("/css/{filename}", cacheHandler("", gzip.GzipHandler(fileHandler)))
	sub.HandleFunc("/css/{timestamp}/{filename}", cacheHandler("/css/", gzip.GzipHandler(fileHandler)))
	sub.HandleFunc("/fonts/{filename}", cacheHandler("", fileHandler))
	sub.HandleFunc("/fonts/{timestamp}/{filename}", cacheHandler("/fonts/", fileHandler))
	sub.HandleFunc("/images/{filename}", cacheHandler("", fileHandler))
	sub.HandleFunc("/images/{timestamp}/{filename}", cacheHandler("/images/", fileHandler))

	r.HandleFunc("/api/v1/users/{login}/sessions", apiNewSessionHandler(s.system, s.sessions)).Methods("POST")
	r.HandleFunc("/logout", logoutHandler(s.system, s.rootPath))
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
		http.ServeFile(w, r, rootPath+"/dist/index.html")
	}
}

func logoutHandler(sys *gohome.System, rootPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sys.Services.EvtBus.Enqueue(&gohome.UserLogoutEvt{
			//TODO: Get the logout details
			Login: "",
		})

		http.ServeFile(w, r, rootPath+"/dist/logout.html")
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
		for _, u := range sys.Users() {
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

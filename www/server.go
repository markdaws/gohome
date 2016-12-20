package www

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
	"github.com/urfave/negroni"
)

type Server struct {
	rootPath       string
	system         *gohome.System
	systemSavePath string
	sessions       *gohome.Sessions
	cfg            *gohome.Config
}

// ListenAndServe creates a new WWW server, that handles API calls and also
// runs the gohome website
func ListenAndServe(
	rootPath string,
	addr string,
	system *gohome.System,
	systemSavePath string,
	sessions *gohome.Sessions,
	cfg *gohome.Config) error {
	server := &Server{
		rootPath:       rootPath,
		system:         system,
		systemSavePath: systemSavePath,
		sessions:       sessions,
		cfg:            cfg,
	}
	return server.listenAndServe(addr)
}

var cacheMutex sync.RWMutex
var cachedFiles = make(map[string][]byte)

// For low horsepower devices such as the Raspberry PI, IO such as reading files and CPU intensive operations
// like GZIP can take multiple seconds for larger files, so here we read the files in and gzip them, the bytes
// are then cached in memory so subsequent requests are faster. Only the first cold hit will be slower to load.
func cacheHandler(prefix string, gzipFile bool, distRoot string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		callStartTime := time.Now().UnixNano()

		originalPath := r.URL.Path

		if prefix != "" {
			// We want to remove the timestamp, that comes before the file name
			temp := strings.Replace(r.URL.Path, prefix, "", -1)
			r.URL.Path = prefix + temp[strings.Index(temp, "/")+1:]
		}

		// Make sure caller is not trying to get out of the base path
		fullPath := distRoot + r.URL.Path
		cleanPath := filepath.Clean(fullPath)
		if strings.HasPrefix(cleanPath, "../") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(r.URL.Path)

		// Since all resources include cache busting values in the URL each time we have a new build,
		// we just set the max cache values here
		w.Header().Add("Cache-Control", "public; max-age=31536000")
		w.Header().Add("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
		w.Header().Add("Content-Type", mime.TypeByExtension(ext))

		if gzipFile {
			w.Header().Add("Content-Encoding", "gzip")
		}

		cacheMutex.RLock()
		b, inCache := cachedFiles[originalPath]
		cacheMutex.RUnlock()

		var readStartTime int64
		var readEndTime int64
		var zipStartTime int64
		var zipEndTime int64
		var acceptsGZIP bool

		if !inCache {
			readStartTime = time.Now().UnixNano()
			var err error
			b, err = ioutil.ReadFile(cleanPath)
			readEndTime = time.Now().UnixNano()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			accept := r.Header.Get("Accept-Encoding")
			if strings.Index(accept, "gzip") != -1 {
				acceptsGZIP = true
			}

			//TODO: Need to cache gzip and unzipped bytes, incase cache
			//gzip bytes but then a client which does not support gzip
			//makes a request
			if gzipFile && acceptsGZIP {
				var gb bytes.Buffer
				zipStartTime = time.Now().UnixNano()
				gz := gzip.NewWriter(&gb)
				if _, err := gz.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if err := gz.Flush(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if err := gz.Close(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				zipEndTime = time.Now().UnixNano()

				b = gb.Bytes()
			}

			cacheMutex.Lock()
			cachedFiles[originalPath] = b
			cacheMutex.Unlock()
		}

		writeStartTime := time.Now().UnixNano()
		w.Write(b)
		writeEndTime := time.Now().UnixNano()

		log.V("webserver - [%s], %dKB, accept gzip: %t, in cache: %t, read:%dms, zip:%dms, write:%dms, total:%dms",
			originalPath, len(b)/1024, acceptsGZIP, inCache, (readEndTime-readStartTime)/1000000, (zipEndTime-zipStartTime)/1000000,
			(writeEndTime-writeStartTime)/1000000, (writeEndTime-callStartTime)/1000000)
	}
}

func (s *Server) listenAndServe(addr string) error {

	r := mux.NewRouter()

	mime.AddExtensionType(".jsx", "text/jsx")
	mime.AddExtensionType(".woff", "application/font-woff")
	mime.AddExtensionType(".woff2", "application/font-woff2")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".ttf", "application/font-sfnt")

	sub := r.PathPrefix("/").Subrouter()

	// For each type of asset, they can be accessed via a direct url like /images/foo.png or they can
	// be accessed with some cache busting value prepended before the filename e.g.
	// /images/12345/foo.png which will redirect to foo.png on the filesystem. This allows you to either
	// put a hash value in the file name of an asset or some cache busting value like the build time in the
	// URL instead of having to rename files
	distPath := s.rootPath + "/dist"
	sub.HandleFunc("/js/{filename}", cacheHandler("", true, distPath))
	sub.HandleFunc("/js/{timestamp}/{filename}", cacheHandler("/js/", true, distPath))
	sub.HandleFunc("/css/{filename}", cacheHandler("", true, distPath))
	sub.HandleFunc("/css/{timestamp}/{filename}", cacheHandler("/css/", true, distPath))
	sub.HandleFunc("/fonts/{filename}", cacheHandler("", false, distPath))
	sub.HandleFunc("/fonts/{timestamp}/{filename}", cacheHandler("/fonts/", false, distPath))
	sub.HandleFunc("/images/{filename}", cacheHandler("", false, distPath))
	sub.HandleFunc("/images/{timestamp}/{filename}", cacheHandler("/images/", false, distPath))

	r.HandleFunc("/api/v1/users/{login}/sessions", apiNewSessionHandler(s.system, s.sessions)).Methods("POST")
	r.HandleFunc("/logout", logoutHandler(s.system, s.rootPath))
	r.HandleFunc("/config", configHandler(s.cfg, s.sessions))
	r.HandleFunc("/system", systemHandler(s.cfg, s.sessions))

	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)
	RegisterSceneHandlers(apiRouter, s)
	RegisterDeviceHandlers(apiRouter, s)
	RegisterDiscoveryHandlers(apiRouter, s)
	RegisterMonitorHandlers(apiRouter, s)
	RegisterAutomationHandlers(apiRouter, s)

	r.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(CheckValidSession(s.sessions)),
		negroni.Wrap(apiRouter),
	))

	r.HandleFunc("/", rootHandler(s.rootPath))

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

func rootHandler(rootPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rootPath+"/dist/index.html")
	}
}

func configHandler(cfg *gohome.Config, sessions *gohome.Sessions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := r.Cookie("sid")
		if err != nil {
			w.Write([]byte("Must be logged in to see this file"))
			return
		}

		_, ok := sessions.Get(sid.Value)
		if !ok {
			w.Write([]byte("Must be logged in to see this file"))
			return
		}

		b, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(b))
	}
}

func systemHandler(cfg *gohome.Config, sessions *gohome.Sessions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := r.Cookie("sid")
		if err != nil {
			w.Write([]byte("Must be logged in to see this file"))
			return
		}

		_, ok := sessions.Get(sid.Value)
		if !ok {
			w.Write([]byte("Must be logged in to see this file"))
			return
		}

		b, err := ioutil.ReadFile(cfg.SystemPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
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

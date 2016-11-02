package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markdaws/gohome"
	errExt "github.com/pkg/errors"
)

func RegisterUserHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/users/{login}/sessions", apiNewSessionHandler(s.system, nil)).Methods("POST")
}

func apiNewSessionHandler(sys *gohome.System, sessions *sessions.CookieStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			respBadRequest("failed to read request body", w)
			return
		}

		login := mux.Vars(r)["login"]
		var user *gohome.User
		for _, u := range sys.Users {
			if u.Login == login {
				user = u
				break
			}
		}
		if user == nil {
			respBadRequest("invalid login: "+login, w)
			return
		}

		var x struct {
			Password string `json:"password"`
		}
		if err = json.Unmarshal(body, &x); err != nil {
			respBadRequest("failed to read request body, invalid JSON", w)
			return
		}

		err = user.VerifyPassword(x.Password)
		if err != nil {
			respBadRequest("invalid password", w)
			return
		}

		session, err := sessions.Get(r, "sid")
		if err != nil {
			respErr(errExt.Wrap(err, "failed to create session"), w)
			return
		}
		session.Values[login] = true
		session.Save(r, w)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			SessionID string `json:"sid"`
		}{
			SessionID: "TODO:",
		})
	}
}

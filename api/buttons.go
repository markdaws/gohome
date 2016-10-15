package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/markdaws/gohome"
)

func apiButtonsHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		buttons := make(buttons, len(system.Buttons), len(system.Buttons))
		var i int32
		for _, button := range system.Buttons {
			buttons[i] = jsonButton{
				ID:       button.ID,
				Name:     button.Name,
				FullName: button.Device.Name + " / " + button.Name,
			}
			i++
		}
		sort.Sort(buttons)
		if err := json.NewEncoder(w).Encode(buttons); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

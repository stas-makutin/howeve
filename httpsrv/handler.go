package httpsrv

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/stas-makutin/howeve/events/handlers"
)

func handleRestart(w http.ResponseWriter, r *http.Request) {
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	handlers.Dispatcher.RequestResponse(r.Context(), &handlers.ConfigGet{RequestHeader: *handlers.NewRequestHeader("")}, reflect.TypeOf(&handlers.ConfigGetResult{}), func(event interface{}) {
		if query := queryFromEvent(event); query != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			err := json.NewEncoder(w).Encode(query)
			if err != nil {
				// TODO
			}
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
}

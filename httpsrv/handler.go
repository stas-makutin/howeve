package httpsrv

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/stas-makutin/howeve/events/handlers"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	handlers.Dispatcher.RequestResponse(r.Context(), &handlers.ConfigGet{RequestHeader: *handlers.NewRequestHeader("")}, reflect.TypeOf(&handlers.ConfigGetResult{}), func(event interface{}) {
		if query := queryFromEvent(event); query != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(query)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
}

package httpsrv

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/stas-makutin/howeve/events/handlers"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	handlers.Dispatcher.RequestResponse(r.Context(), &handlers.ConfigGet{}, reflect.TypeOf(&handlers.ConfigData{}), func(event interface{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(event.(*handlers.ConfigData).Config)
	})
}

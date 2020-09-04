package httpsrv

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/stas-makutin/howeve/eventh"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	eventh.Dispatcher.RequestResponse(r.Context(), &eventh.ConfigGet{}, reflect.TypeOf(&eventh.ConfigData{}), func(event interface{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(event.(*eventh.ConfigData).Config)
	})
}

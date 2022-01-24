package httpsrv

import (
	"net/http"
	"strings"

	"github.com/stas-makutin/howeve/config"
)

type asset config.HTTPAsset

func (a *asset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, a.Route)
	w.Write([]byte(path))
}

package httpsrv

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stas-makutin/howeve/config"
)

type asset config.HTTPAsset

func (a *asset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, a.Route) {
		http.NotFound(w, r)
		return
	}

	path := r.URL.Path[len(a.Route):]
	root := false
	if path == "" {
		path = a.Path
		root = true
	} else {
		path = filepath.Join(a.Path, path)
	}
	path = filepath.Clean(path)

	if len(a.Excludes) > 0 {
		for _, pattern := range a.Excludes {
			if m, err := filepath.Match(pattern, path); m {
				http.NotFound(w, r)
				return
			} else if err != nil {
				appendLogFields(r, fmt.Sprintf("invalid exclude pattern %s: %v", pattern, err.Error()))
			}
			if m, _ := filepath.Match(pattern, filepath.Base(path)); m {
				http.NotFound(w, r)
				return
			}
		}
	}
	if len(a.Includes) > 0 {
		for _, pattern := range a.Includes {
			if m, err := filepath.Match(pattern, path); !m {
				if m, _ := filepath.Match(pattern, filepath.Base(path)); !m {
					http.NotFound(w, r)
					return
				}
			} else if err != nil {
				appendLogFields(r, fmt.Sprintf("invalid include pattern %s: %v", pattern, err.Error()))
			}
		}
	}

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			appendLogFields(r, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if fi.IsDir() {
		noIndexFile := true
		if root && len(a.IndexFiles) > 0 {
			for _, indexFile := range a.IndexFiles {
				indexPath := filepath.Join(path, indexFile)
				if ifi, err := os.Stat(indexPath); err == nil && !ifi.IsDir() {
					appendLogFields(r, indexFile)
					fi = ifi
					path = indexPath
					noIndexFile = false
					break
				}
			}
		}
		if noIndexFile {
			if (a.Flags & config.HAFDirListing) != 0 {
				http.ServeFile(w, r, path)
			} else {
				http.NotFound(w, r)
			}
			return
		}
	}

	f, err := os.Open(path)
	if err != nil {
		appendLogFields(r, err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	http.ServeContent(w, r, path, fi.ModTime(), f)
}

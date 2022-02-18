package httpsrv

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/log"
)

// log constants
const (
	// operation
	haOpAddFromConfig = "C"

	// operation codes
	haOcRouteConflict = "R"
	// exclusion pattern invalid
	haOcExcludeInvalid = "E"
	// inclusion pattern invalid
	haOcIncludeInvalid = "I"
)

type asset config.HTTPAsset

func (a *asset) valid(routes map[string]struct{}) bool {
	rc := true
	if _, ok := routes[a.Route]; ok || a.Route == "" {
		log.Report(log.SrcHTTPAssets, haOpAddFromConfig, haOcRouteConflict, a.Route)
		rc = false
	}
	for _, pattern := range a.Excludes {
		if _, err := filepath.Match(pattern, "/"); err != nil {
			log.Report(log.SrcHTTPAssets, haOpAddFromConfig, haOcExcludeInvalid, a.Route, pattern)
			rc = false
		}
	}
	for _, pattern := range a.Includes {
		if _, err := filepath.Match(pattern, "/"); err != nil {
			log.Report(log.SrcHTTPAssets, haOpAddFromConfig, haOcIncludeInvalid, a.Route, pattern)
			rc = false
		}
	}
	return rc
}

func (a *asset) checkVisibility(path string) bool {
	name := filepath.Base(path)
	if name == "" {
		return false
	}
	if (a.Flags & config.HAFShowHidden) == 0 {
		if name[0:1] == "." {
			return false
		}
	}
	if len(a.Excludes) > 0 {
		for _, pattern := range a.Excludes {
			if m, err := filepath.Match(pattern, path); m || err != nil {
				return false
			}
			if m, err := filepath.Match(pattern, filepath.Base(path)); m || err != nil {
				return false
			}
		}
	}
	if len(a.Includes) > 0 {
		for _, pattern := range a.Includes {
			if m, err := filepath.Match(pattern, path); !m || err != nil {
				if m, err := filepath.Match(pattern, filepath.Base(path)); !m || err != nil {
					return false
				}
			}
		}
	}
	return true
}

func (a *asset) dirListing(w http.ResponseWriter, r *http.Request, path string, modtime time.Time, root bool) {
	if !modtime.IsZero() {
		w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
		if r.Method == "GET" || r.Method == "HEAD" {
			ims := r.Header.Get("If-Modified-Since")
			for _, layout := range []string{http.TimeFormat, time.RFC850, time.ANSIC} {
				t, err := time.Parse(layout, ims)
				if err == nil {
					mt := modtime.Truncate(time.Second)
					if mt.Before(t) || mt.Equal(t) {
						return
					}
					break
				}
			}
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	if !root {
		fmt.Fprintf(w, "<a href=\"..\">..</a>\n")
	}
	for _, file := range files {
		name := file.Name()
		if a.checkVisibility(filepath.Join(path, name)) {
			if file.IsDir() {
				name += "/"
			}
			url := url.URL{Path: name}
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", html.EscapeString(url.String()), html.EscapeString(name))
		}
	}
	fmt.Fprintf(w, "</pre>\n")
}

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

	if !a.checkVisibility(path) {
		http.NotFound(w, r)
		return
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) && (a.Flags&config.HAFFlat) != 0 {
		path = filepath.Clean(filepath.Join(a.Path, filepath.Base(path)))
		fi, err = os.Stat(path)
		if os.IsNotExist(err) {
			path = filepath.Clean(a.Path)
			root = true
			fi, err = os.Stat(path)
		}
	}
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
				a.dirListing(w, r, path, fi.ModTime(), root)
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

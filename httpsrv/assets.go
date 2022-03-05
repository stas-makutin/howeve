package httpsrv

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/stas-makutin/howeve/api"
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

type asset api.HTTPAsset

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

func (a *asset) checkVisibility(trgPath string) bool {
	name := filepath.Base(trgPath)
	if name == "" {
		return false
	}
	if (a.Flags & api.HAFShowHidden) == 0 {
		if name[0:1] == "." {
			return false
		}
	}
	if len(a.Excludes) > 0 {
		for _, pattern := range a.Excludes {
			if m, err := filepath.Match(pattern, trgPath); m || err != nil {
				return false
			}
			if m, err := filepath.Match(pattern, filepath.Base(trgPath)); m || err != nil {
				return false
			}
		}
	}
	if len(a.Includes) > 0 {
		for _, pattern := range a.Includes {
			if m, err := filepath.Match(pattern, trgPath); !m || err != nil {
				if m, err := filepath.Match(pattern, filepath.Base(trgPath)); !m || err != nil {
					return false
				}
			}
		}
	}
	return true
}

func (a *asset) dirListing(w http.ResponseWriter, r *http.Request, trgPath string, modtime time.Time, root bool) {
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

	files, err := os.ReadDir(trgPath)
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
		if a.checkVisibility(filepath.Join(trgPath, name)) {
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

	relPath := r.URL.Path[len(a.Route):]
	trgPath := relPath
	root := false
	if trgPath == "" {
		trgPath = a.Path
		root = true
	} else {
		trgPath = filepath.Join(a.Path, trgPath)
	}
	trgPath = filepath.Clean(trgPath)

	if !a.checkVisibility(trgPath) {
		http.NotFound(w, r)
		return
	}

	fi, err := os.Stat(trgPath)
	if os.IsNotExist(err) && (a.Flags&api.HAFFlat) != 0 {
		// in flat mode all paths "flats" into asset directory to support client-based routing
		// main idea of this mode that if some path is not exists it always translated to asset path
		// example:
		//   asset route: /a/
		//   asset path: my/page
		//   assets which exists: my/page/file, my/page/dir/file1, my/page/dir/sub/file2
		//   request paths translation:
		//      /a/ -> my/page
		//      /a/file -> my/page/file
		//      /a/client_route -> my/page
		//      /a/dir/file1 -> my/page/dir/file1
		//      /a/dir/sub/file2 -> my/page/dir/file2
		// and special cases:
		//      /a/b/b/b/file -> my/page/file
		//      /a/b/b/b/dir/file1 -> my/page/dir/file1
		//      /a/b/b/b/dir/sub/file1 -> my/page/dir/sub/file1
		tryPath := ""
		relPath, base := path.Split(relPath)
		count := 0
		for {
			relPath = path.Dir(relPath)
			if base != "" {
				tryPath = path.Join(base, tryPath)
				trgPath = filepath.Clean(filepath.Join(a.Path, tryPath))
				fi, err = os.Stat(trgPath)
				if err == nil {
					break
				}
				if count > 64 {
					// limit the number of allowed subdirectories for protection
					http.NotFound(w, r)
					return
				}
				count++
			}
			if relPath == "." || relPath == "" {
				break
			}
			base = path.Base(relPath)
		}
		if err != nil {
			trgPath = filepath.Clean(a.Path)
			root = true
			fi, err = os.Stat(trgPath)
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
				indexPath := filepath.Join(trgPath, indexFile)
				if ifi, err := os.Stat(indexPath); err == nil && !ifi.IsDir() {
					appendLogFields(r, indexFile)
					fi = ifi
					trgPath = indexPath
					noIndexFile = false
					break
				}
			}
		}
		if noIndexFile {
			if (a.Flags & api.HAFDirListing) != 0 {
				a.dirListing(w, r, trgPath, fi.ModTime(), root)
			} else {
				http.NotFound(w, r)
			}
			return
		}
	}

	f, err := os.Open(trgPath)
	if err != nil {
		appendLogFields(r, err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	http.ServeContent(w, r, trgPath, fi.ModTime(), f)
}

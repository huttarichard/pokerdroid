package studio

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	app "github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/studio/internal/fsutil"
	"github.com/pokerdroid/poker/studio/internal/jsc"
)

func NewHttpHandler(fs fs.FS) (http.HandlerFunc, error) {
	fs, err := fsutil.NewTemplate(fs, []string{"index.html"}, func(f string) (any, error) {
		var data struct {
			Data       template.JS
			Entrypoint string
			Name       string
			Version    string
			Loaded     bool
		}

		data.Name = app.Name
		data.Version = app.GitCommit

		meta, err := jsc.ReadMetafile(fs, "dist/meta.json")
		if err != nil {
			return data, nil
		}

		file, ok := meta.Outputs.FindByEntrypoint("src/main.tsx")
		if !ok {
			return data, nil
		}

		data.Entrypoint = file
		data.Loaded = true
		data.Data = template.JS([]byte("{}"))
		return data, nil
	})

	if err != nil {
		return nil, err
	}

	handler := NewSpaHandlerFunc(SpaHandlerParams{
		FS:         fs,
		PublicDirs: []string{"public", "dist"},
		Index:      "index.html",
	})

	return handler, nil
}

type SpaHandlerParams struct {
	FS              fs.FS
	PublicDirs      []string
	Index           string
	NotFoundHandler http.HandlerFunc
}

type SpaHandler struct {
	SpaHandlerParams
}

// NewSpaHandler creates a new handler that serves a SPA.
//
// By default it servers Index file for all routes except for the
// ones that start with DistDir or PublicDir.
//
// PublicDir directory is used to serve static files as in root.
func NewSpaHandler(p SpaHandlerParams) (h http.Handler) {
	p.Index = CleanURLPath(p.Index)
	if p.Index == "" {
		p.Index = "/index.html"
	}

	for i, d := range p.PublicDirs {
		p.PublicDirs[i] = CleanURLPath(d)
	}

	asset := &SpaHandler{
		SpaHandlerParams: p,
	}

	return asset
}

func NewSpaHandlerFunc(p SpaHandlerParams) http.HandlerFunc {
	h := NewSpaHandler(p)
	return h.ServeHTTP
}

func (h *SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the path is clean to prevent directory traversal vulnerabilities
	upath := CleanURLPath(r.URL.Path)
	r.URL.Path = upath

	fsh := http.FS(h.FS)

	switch {
	case upath == "/":
	case upath == h.Index:
		h.ServeIndex(w)
		return
	default:
	}

	var found http.File
	for _, p := range h.PublicDirs {
		p = path.Join(p, upath)
		f, err := fsh.Open(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		found = f
		r.URL.Path = p
		break
	}

	if found == nil {
		h.ServeNotFound(w, r)
		return
	}

	d, err := found.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if d.IsDir() {
		h.ServeNotFound(w, r)
		return
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), found)
}

func (h *SpaHandler) ServeNotFound(w http.ResponseWriter, r *http.Request) {
	if h.NotFoundHandler != nil {
		h.NotFoundHandler(w, r)
		return
	}
	h.ServeIndex(w)
}

func (h *SpaHandler) ServeIndex(w http.ResponseWriter) {
	p := strings.TrimLeft(h.Index, "/")
	f, err := h.FS.Open(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	n, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hd := w.Header()
	hd.Set("Content-Type", "text/html")
	hd.Set("Content-Length", fmt.Sprint(n.Size()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func CleanURLPath(p string) string {
	p = path.Clean(p)
	p = strings.TrimLeft(p, "/")
	return fmt.Sprintf("/%s", p)
}

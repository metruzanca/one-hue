// Package serve hosts the theme preview web app: an HTMX + CodeMirror page
// that previews, creates and edits theme variants, writes themes.toml and
// regenerates the editor theme files on save.
package serve

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/metruzanca/one-hue-theme/internal/gen"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

//go:embed web/index.html
var webFS embed.FS

var tmpl = template.Must(template.ParseFS(webFS, "web/index.html"))

// Options configures the server.
type Options struct {
	// Addr is the listen address, e.g. ":8080".
	Addr string
	// ConfigPath is the themes TOML file to read and write.
	ConfigPath string
	// OutDir is the generated themes output directory.
	OutDir string
}

// server holds the shared state of the running app.
type server struct {
	opts Options

	// mu guards config reads/writes and the rebuild so concurrent requests
	// cannot interleave a save with a build.
	mu sync.Mutex
}

// Serve starts the HTTP server and blocks until it stops.
func Serve(opts Options) error {
	if opts.Addr == "" {
		opts.Addr = ":8080"
	}
	s := &server{opts: opts}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleIndex)
	mux.HandleFunc("GET /api/themes", s.handleList)
	mux.HandleFunc("POST /api/preview", s.handlePreview)
	mux.HandleFunc("PUT /api/themes", s.handleSave)

	log.Printf("one-hue serve listening on http://localhost%s", opts.Addr)
	return http.ListenAndServe(opts.Addr, mux)
}

// handleIndex renders the full app page with the theme list.
func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	configs, err := theme.Load(s.opts.ConfigPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "index.html", configs); err != nil {
		log.Printf("render index: %v", err)
	}
}

// handleList returns the configured themes as JSON.
func (s *server) handleList(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	configs, err := theme.Load(s.opts.ConfigPath)
	s.mu.Unlock()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	out := make([]map[string]string, 0, len(configs))
	for _, c := range configs {
		out = append(out, map[string]string{"slug": c.Slug, "name": c.Name, "accent": c.Accent})
	}
	writeJSON(w, http.StatusOK, out)
}

// previewRequest is the body of POST /api/preview.
type previewRequest struct {
	Name   string `json:"name"`
	Accent string `json:"accent"`
}

// handlePreview derives a theme from a name/accent and returns the preview
// payload. It writes nothing.
func (s *server) handlePreview(w http.ResponseWriter, r *http.Request) {
	var req previewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	name := req.Name
	if name == "" {
		name = "Preview"
	}
	cfg, err := theme.ConfigFromVariant(theme.Variant{Name: name, Accent: req.Accent})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, theme.PreviewOf(theme.Build(cfg)))
}

// saveRequest is the body of PUT /api/themes.
type saveRequest struct {
	// OriginalSlug identifies the theme being edited; empty creates a new one.
	OriginalSlug string `json:"originalSlug"`
	Name         string `json:"name"`
	Accent       string `json:"accent"`
}

// handleSave upserts a theme into themes.toml and rebuilds every target,
// replying with the refreshed list fragment and a flash message.
func (s *server) handleSave(w http.ResponseWriter, r *http.Request) {
	var req saveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeList(w, http.StatusOK, "error", "invalid JSON body")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	variants, err := theme.LoadVariants(s.opts.ConfigPath)
	if err != nil {
		s.writeList(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	cfg, err := theme.ConfigFromVariant(theme.Variant{Name: req.Name, Accent: req.Accent})
	if err != nil {
		s.writeList(w, http.StatusOK, "error", err.Error())
		return
	}

	idx := -1
	if req.OriginalSlug != "" {
		for i, v := range variants {
			if theme.Slugify(v.Name) == req.OriginalSlug {
				idx = i
				break
			}
		}
	}
	for i, v := range variants {
		if i != idx && theme.Slugify(v.Name) == cfg.Slug {
			s.writeList(w, http.StatusOK, "error",
				fmt.Sprintf("a theme named %q already exists", req.Name))
			return
		}
	}

	entry := theme.Variant{Name: cfg.Name, Accent: cfg.Accent}
	if idx >= 0 {
		variants[idx] = entry
	} else {
		variants = append(variants, entry)
	}

	if err := theme.SaveVariants(s.opts.ConfigPath, variants); err != nil {
		s.writeList(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	if err := rebuild(s.opts.ConfigPath, s.opts.OutDir); err != nil {
		s.writeList(w, http.StatusInternalServerError, "error", "saved, but rebuild failed: "+err.Error())
		return
	}

	configs, err := theme.Load(s.opts.ConfigPath)
	if err != nil {
		s.writeList(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	s.writeListConfigs(w, http.StatusOK, configs, "ok", fmt.Sprintf("Saved %s and rebuilt themes", cfg.Name))
}

// writeList reloads the configs and renders a list fragment, used for error
// paths that still need to return the current list.
func (s *server) writeList(w http.ResponseWriter, status int, kind, msg string) {
	configs, err := theme.Load(s.opts.ConfigPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeListConfigs(w, status, configs, kind, msg)
}

// writeListConfigs renders the list fragment plus an out-of-band flash.
func (s *server) writeListConfigs(w http.ResponseWriter, status int, configs []theme.Config, kind, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "themes-list", configs); err != nil {
		log.Printf("render list: %v", err)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "flash", map[string]string{"Kind": kind, "Message": msg}); err != nil {
		log.Printf("render flash: %v", err)
	}
}

// rebuild regenerates every editor target from the config file.
func rebuild(configPath, outDir string) error {
	configs, err := theme.Load(configPath)
	if err != nil {
		return err
	}
	built := make([]*theme.Built, len(configs))
	for i := range configs {
		built[i] = theme.Build(configs[i])
	}
	files, err := gen.RenderAll(built)
	if err != nil {
		return err
	}
	return gen.WriteFiles(outDir, files)
}

// writeJSON writes v as a JSON response with the given status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}

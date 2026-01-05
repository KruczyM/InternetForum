package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"forum/ui"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"

)

type Handler struct {
	DB             *sql.DB
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	TemplateCache  map[string]*template.Template
}

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	//its giveing whole stack trace so we know in which line in code it happening
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.ErrorLog.Output(2, trace)

	w.WriteHeader(http.StatusInternalServerError)
	h.render(w, http.StatusInternalServerError, "error.html", nil)
}

// just to make work a little bit easier and faster
func (h *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(w,r)
	h.render(w, http.StatusNotFound, "404.html", data)
}

// Return true if the current request is from an authenticated user, otherwise return false.
func (h *Handler) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (h *Handler) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.TemplateCache[page]
	if !ok {
		h.serverError(w, fmt.Errorf("template %s does not exist", page))
		return
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, data)
	if err != nil {
		h.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	pages, err := fs.Glob(ui.Files, "html/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(funcMap).ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}
		cache[name] = templateSet
	}

	return cache, nil
}

func (h *Handler) test500(w http.ResponseWriter, r *http.Request) {
	panic("test 500 error")
}

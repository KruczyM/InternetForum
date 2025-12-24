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

	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	DB             *sql.DB
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	SessionManager *scs.SessionManager
	TemplateCache  map[string]*template.Template
}

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	//its giveing whole stack trace so we know in which line in code it happening
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// just to make work a little bit easier and faster
func (h *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (h *Handler) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}

// Return true if the current request is from an authenticated user, otherwise return false.
func (h *Handler) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

type templateData struct {
	Form            any
	Flash           string
	IsAuthenticated bool
	//CSRFToken       string
}

func (h *Handler) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		Flash:           h.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: h.isAuthenticated(r),
		// CSRFToken:       nosurf.Token(r),
	}
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
	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded filesystem which match the pattern 'html/*.tmpl'. This essentially
	// gives us a slice of all the 'page' templates for the application
	pages, err := fs.Glob(ui.Files, "html/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// it will give us the name of the file without the extension
		name := filepath.Base(page)
		// Use template.New() to create a new template with the given name. This will return an error if the template name is invalid.
		templateSet, err := template.New(name).ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}

package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"database/sql"
	"log"
	"github.com/alexedwards/scs/v2"

)

type Handler struct {
    DB *sql.DB
    InfoLog *log.Logger
    ErrorLog *log.Logger
	SessionManager *scs.SessionManager
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

// // Return true if the current request is from an authenticated user, otherwise return false.
// func (h *Handler) isAuthenticated(r *http.Request) bool {
// 	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
// 	if !ok {
// 		return false
// 	}
// 	return isAuthenticated
// }
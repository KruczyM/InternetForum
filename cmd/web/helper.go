package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"database/sql"
	"log"

	"github.com/alexedwards/scs/v2"
)

type application struct {
    errorLog *log.Logger
    infoLog  *log.Logger

    db *sql.DB
	sessionManager *scs.SessionManager
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
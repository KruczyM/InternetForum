package main

import (
	"forum/cmd/web/handlers"
	"forum/internal/db"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)


func main() {

infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
db, err := database.InitDB("data/forum.db")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

// create a new session manager with information about the cookies and the store
sessionManager := scs.New()
sessionManager.Lifetime = 24 * time.Hour
sessionManager.IdleTimeout = 30 * time.Minute
sessionManager.Cookie.HttpOnly = true
sessionManager.Cookie.SameSite = http.SameSiteLaxMode

// localy false but when we deploy it to the server we should set it to true 
sessionManager.Cookie.Secure = false

// it will manage the sessions and store them in the database
sessionManager.Store = sqlite3store.New(db)

handler := &handlers.Handler{
	DB:       db,
	InfoLog:  infoLog,
	ErrorLog: errorLog,
	SessionManager: sessionManager,
}
	

	appRouter := handler.Routes()
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", sessionManager.LoadAndSave(appRouter)))




}


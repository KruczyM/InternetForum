package main

import (
	"encoding/gob"
	"forum/cmd/web/handlers"
	database "forum/internal/db"
	"log"
	"net/http"
	"os"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	gob.Register(&handlers.FlashMessage{})

	db, err := database.InitDB("data/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()




	templateCache, err := handlers.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	handler := &handlers.Handler{
		DB:             db,
		InfoLog:        infoLog,
		ErrorLog:       errorLog,
		TemplateCache:  templateCache,
	}

	appRouter := handler.Routes()
	log.Println("Server starting on http://localhost:8080")

	err = http.ListenAndServe(":8080", appRouter)
	if err != nil {
	log.Fatal(err)
}

}

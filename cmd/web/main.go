package main

import (
	"forum/internal/db"
	"forum/cmd/web/handlers"
	"log"
	"net/http"
	"os"
)


func main() {

infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
db, err := database.InitDB("localData/forum.db")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

handler := &handlers.Handler{
	DB:       db,
	InfoLog:  infoLog,
	ErrorLog: errorLog,
}

	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/post/create", handler.CreatePost)
	http.HandleFunc("/auth/register", handler.Register)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

// app := &application{
// 	errorLog: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
// 	infoLog:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
// 	db: db,
// }


}


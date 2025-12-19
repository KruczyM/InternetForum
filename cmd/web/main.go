package main

import (
	"forum/internal/db"
	"log"
)


func main() {

	db, err := database.Open("localData/forum.db")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

if err := database.RunMigrations(db); err != nil {
	log.Fatal(err)
}

seedUsers(db)

// app := &application{
// 	errorLog: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
// 	infoLog:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
// 	db: db,
// }


}


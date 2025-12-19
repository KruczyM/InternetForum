package main

import (
	"forum/internal/db"
	"log"
)


func main() {
	db, err := database.Open("./data/forum.db")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

if err := database.RunMigrations(db); err != nil {
	log.Fatal(err)
}
}

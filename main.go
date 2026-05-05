package main

import (
	"log"
	"net/http"
	"poker-engine/db"
	"poker-engine/routes"
)

func main() {
	config := GetConfig()

	db, err := db.DBConnection(config.DatabaseURL)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = http.ListenAndServe(":8080", routes.Router(db))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

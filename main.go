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

	router := routes.NewAppRouter(db)

	log.Printf("Listening on port: 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

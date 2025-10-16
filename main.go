package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}
	webDir := os.Getenv("WEB_DIR")
	if webDir == "" {
		webDir = "./web"
	}
	dbName := os.Getenv("TODO_DBFILE")
	if dbName == "" {
		dbName = "scheduler.db"
	}

	srv := server.Run(webDir)

	err = db.Init("scheduler.db")
	if err != nil {
		log.Fatal("error init db")
	}

	err = http.ListenAndServe(port, srv)
	if err != nil {
		log.Fatal("error starting server")
	}

}

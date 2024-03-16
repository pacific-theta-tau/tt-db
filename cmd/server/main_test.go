package main

import (
	"fmt"

	"github.com/pacific-theta-tau/tt-db/api"
	"github.com/pacific-theta-tau/tt-db/db"
)

func runTestServer() {
	config := api.Config{
		Port:        "3000",
		DatabaseURL: ":memory:",
	}
	db := db.NewPostgresDB()
	app := api.NewApplication(config, db)
	fmt.Println("Starting server on port", app.Config.Port, "...")
	app.Serve()
	// TODO: create seeds and migration for test DB
}

func Test
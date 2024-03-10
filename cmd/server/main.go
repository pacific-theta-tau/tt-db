// Main driver. All application and database setups happens here
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pacific-theta-tau/tt-db/api"
	"github.com/pacific-theta-tau/tt-db/db"
)

func main() {
	// Make sure to setup .env file with all necessary configs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	config := api.Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	db := db.NewPostgresDB()
	app := api.NewApplication(config, db)
	fmt.Println("Starting server on port", app.Config.Port, "...")
	app.Serve()
}

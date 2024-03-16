// Main driver. All application and database setups happens here
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pacific-theta-tau/tt-db/api"
	"github.com/pacific-theta-tau/tt-db/db"
)

func main() {
	// Make sure to setup .env file with all necessary configs
	devFlag := flag.String("env", "dev", "Environment to serve app. Options: dev (default) | prod")
	flag.Parse()
	var port string
	var databaseUrl string

	if *devFlag == "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file.")
		}
	}

	port = os.Getenv("PORT")
	databaseUrl = os.Getenv("DATABASE_URL")

	fmt.Println("Starting app in", *devFlag, "environment.")
	config := api.Config{
		Port:        port,
		DatabaseURL: databaseUrl,
	}
	fmt.Println("Using config:", config)
	db := db.NewPostgresDB()
	app := api.NewApplication(config, db)
	fmt.Println("Serving app on port", app.Config.Port, "...")
	app.Serve()
}

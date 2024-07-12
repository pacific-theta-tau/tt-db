// Main driver. All application and database setups happens here
package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pacific-theta-tau/tt-db/api"
	"github.com/pacific-theta-tau/tt-db/db"
)

func main() {
	// Loading environment variables
	devFlag := flag.String("env", "dev", "Environment to serve app. Options: dev (default) | prod")
	flag.Parse()
	err := godotenv.Load(*devFlag + ".env")
	if err != nil {
		log.Fatal("Error loading .env file. Make sure to have setup the appropriate .env file")
	}

	log.Println("** Initializing app in", *devFlag, "environment **")
	// Setup app configurations
	app_port := os.Getenv("APP_PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("ERROR: environment variable DatabaseURL not set")
	}

	// Connect to DB and serve API
	db := db.NewPostgresDB(databaseURL)
	app := api.NewApplication(db, app_port)
	log.Printf("Serving app on port %s ...", app.Port)
	app.Serve()
}

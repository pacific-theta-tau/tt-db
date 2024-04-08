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
	// Loading environment variables
	devFlag := flag.String("env", "dev", "Environment to serve app. Options: dev (default) | prod")
	flag.Parse()
	err := godotenv.Load(*devFlag + ".env")
	if err != nil {
		log.Fatal("Error loading .env file. Make sure to have setup the appropriate .env file")
	}

	fmt.Println("** Initializing app in", *devFlag, "environment **")
	// Setup app configurations
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("ERROR: environment variable DatabaseURL not set")
	}

	// Connect to DB and serve API
	db := db.NewPostgresDB(databaseURL)
	app := api.NewApplication(db, port)
	fmt.Println("Serving app on port", app.Port, "...")
	app.Serve()
}

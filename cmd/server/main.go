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
	databaseUrl := os.Getenv("DATABASE_URL")
	config := api.Config{
		Port:        port,
		DatabaseURL: databaseUrl,
	}
	fmt.Println("Using configs:", config)

	// Connect to DB and serve API
	db := db.NewPostgresDB()
	app := api.NewApplication(config, db)
	fmt.Println("Serving app on port", app.Config.Port, "...")
	app.Serve()
}

package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// ConnectPostgresDB estabilish connection with a postgreSQL database defined in the environment variable "DATABASE_URL".
// Returns the database connection as a pgx.Conn object
func ConnectPostgresDB() (*pgx.Conn, error) {
	// Make sure to have setup a .env file with database configs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	// Establish connection to postgres DB
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}

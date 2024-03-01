package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectPostgresDB() (*pgx.Conn, error) {
	// Load DB configs from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	// Connect to postgres DB
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}

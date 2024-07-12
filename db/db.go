package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresDB struct {
	Conn        *sql.DB
	DatabaseURL string
}

// Constructor for PostgresDB struct
func NewPostgresDB(databaseURL string) *PostgresDB {
	return &PostgresDB{
		DatabaseURL: databaseURL,
	}
}

// Estabilishes connection with a postgreSQL database defined in the environment variable "DATABASE_URL".
func (db *PostgresDB) Connect() {
    log.Println(fmt.Sprintf("Connecting to Database URL %s", db.DatabaseURL))

	// Establish connection to postgres DB
	conn, err := sql.Open("pgx", db.DatabaseURL)
	if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
	}

	err = testDB(conn)
	if err != nil {
        log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Connected to Database successfully!")
	db.Conn = conn
}

// Test connection with database by sending a ping
func testDB(conn *sql.DB) error {
	err := conn.Ping()
	if err != nil {
		return err
	}

	return nil
}

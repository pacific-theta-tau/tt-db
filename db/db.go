package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresDB struct {
	Conn *sql.DB
}

// Constructor for PostgresDB struct
func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

// Estabilishes connection with a postgreSQL database defined in the environment variable "DATABASE_URL".
func (db *PostgresDB) Connect(databaseURL string) {
	// Establish connection to postgres DB
	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		fmt.Println("database url:", databaseURL)
		log.Fatal(err)
	}

	err = testDB(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database successfully!")
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

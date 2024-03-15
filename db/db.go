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
func (db *PostgresDB) Connect(dsn string) {
	// Establish connection to postgres DB
	// dsn = "postgresql://Theta%20Tau%20Database_owner:dWotgmw7q1QH@ep-fancy-smoke-a6wr39jt.us-west-2.aws.neon.tech/Theta%20Tau%20Database?sslmode=require"
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
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

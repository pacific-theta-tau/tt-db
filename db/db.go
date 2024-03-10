package db

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
)

type PostgresDB struct {
	db *sql.DB
}

// Constructor for PostgresDB struct
func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

// Connect() estabilishes connection with a postgreSQL database defined in the environment variable "DATABASE_URL".
// Returns the database connection as a pgx.Conn object
func (d *PostgresDB) Connect(addr string) {
	// Make sure to have setup a .env file with database configs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	// // Establish connection to postgres DB
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	conn, err := sql.Open("pgx", addr)
	if err != nil {
		log.Fatal(err)
	}
	d.db = conn
}

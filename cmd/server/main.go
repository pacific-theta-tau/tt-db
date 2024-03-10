package main

import (
	"flag"
	"fmt"

	"github.com/pacific-theta-tau/tt-db/api"
	"github.com/pacific-theta-tau/tt-db/db"
)

func main() {
	port := flag.String("port", ":3000", "Port for server to listen and serve to")
	flag.Parse()

	db := db.NewPostgresDB()

	app := api.NewApplication(*port, db)
	fmt.Println("Server running on port", *port)
	app.Serve()
}

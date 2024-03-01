package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pacific-theta-tau/tt-db/db"
)

func main() {
	dbConn, err := db.ConnectPostgresDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
	}
	defer dbConn.Close(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(":3000", r)
}

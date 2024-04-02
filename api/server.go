// This package defines server and routes
package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pacific-theta-tau/tt-db/api/handlers"
	"github.com/pacific-theta-tau/tt-db/db"
)

type Config struct {
	Port        string
	DatabaseURL string
}

type Application struct {
	Database *db.PostgresDB
	Config   Config
}

func NewApplication(config Config, db *db.PostgresDB) *Application {
	return &Application{
		Config:   config,
		Database: db,
	}
}

// Connect to database, start routers, and serve app
func (app *Application) Serve() {
	// Connect to database
	app.Database.Connect(app.Config.DatabaseURL)

	// Start routers and middleware
	handler := handlers.NewHandler(app.Database.Conn)
	routes := setupRoutes(handler)

	//TODO: cleaner address
	addr := fmt.Sprint(":", app.Config.Port)
	err := http.ListenAndServe(addr, routes)
	if err != nil {
		log.Fatal(err)
	}
}

// Setup mux with all middleware and routes
func setupRoutes(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/api/brothers", handler.GetAllBrothers)
	r.Get("/api/brothers/{rollCall}", handler.GetBrotherByRollCall)
	r.Post("/api/brothers", handler.AddBrother)
	r.Delete("/api/brothers", handler.RemoveBrother)

	return r
}

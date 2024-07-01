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

// TODO: Create a Pgpool instead using dependency injection
type Application struct {
	Database    *db.PostgresDB
	DatabaseURL string
	Port        string
}

// Constructor for Application struct
func NewApplication(db *db.PostgresDB, port string) *Application {
	return &Application{
		Database: db,
		Port:     port,
	}
}

// Connect to database, start routers, and serve app
func (app *Application) Serve() {
    log.Println("-- Application.Serve() --")
	// Connect to database
	app.Database.Connect()

	// Start routers and middleware
	handler := handlers.NewHandler(app.Database.Conn)
	routes := setupRoutes(handler)

	//TODO: cleaner address
	addr := fmt.Sprint(":", app.Port)
	log.Printf("App address: %s", addr)
	err := http.ListenAndServe(addr, routes)
	if err != nil {
        log.Fatalf("Error while serving application: %v", err)
	}
}

// Setup mux with all middleware and routes
func setupRoutes(handler *handlers.Handler) *chi.Mux {
    log.Println("Setting up routes...")
	r := chi.NewRouter()
    // TODO: look into slog for structured logging
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	// brothers endpoint
	r.Get("/api/brothers", handler.GetAllBrothers)
	r.Get("/api/brothers/{rollCall}", handler.GetBrotherByRollCall)
	r.Post("/api/brothers", handler.AddBrother)
	r.Put("/api/brothers", handler.UpdateBrother)
	r.Delete("/api/brothers", handler.RemoveBrother)

    // TODO: events endpoint
	r.Get("/api/events", handler.GetAllEvents)

    // TODO: attendance endpoints

	return r
}

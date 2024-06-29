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
	// Connect to database
	app.Database.Connect()

	// Start routers and middleware
	handler := handlers.NewHandler(app.Database.Conn)
	routes := setupRoutes(handler)

	//TODO: cleaner address
	addr := fmt.Sprint(":", app.Port)
	fmt.Println("address:", addr)
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

	// brothers endpoint
	r.Get("/api/brothers", handler.GetAllBrothers)
	r.Get("/api/brothers/{rollCall}", handler.GetBrotherByRollCall)
	r.Post("/api/brothers", handler.AddBrother)
	r.Put("/api/brothers", handler.UpdateBrother)
	r.Delete("/api/brothers", handler.RemoveBrother)

	// events endpoint
	r.Get("/api/events", handler.GetAllEvents)

	return r
}

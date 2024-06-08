package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to access the API
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},                                                       // allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 // allow only REST requests
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // allow only these headers
		ExposedHeaders:   []string{"Link"},                                                    // allow only these headers to be exposed
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ready")) // add a heartbeat endpoint

	mux.Get("/healthz", app.HandleHealthz) // add a heartbeat endpoint

	mux.Get("/path", app.HandleGetPath)
	//mux.Post("/routes", app.HandleCurrentPosition) // handle subscription to manage current shipment position

	return mux

}

package main

import (
	"net/http"
	"os"
)

// Handles the health endpoint for Dapr
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5280"
	if value, ok := os.LookupEnv("DAPR_HTTP_PORT"); ok {
		daprHttpPort = value
	}
	_, err := http.Get("http://localhost:" + daprHttpPort + "/v1.0/healthz")

	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		os.Exit(1)
	}

	app.writeJSON(w, http.StatusOK, "Healthy")
}

// Handle Get path endpoint
func (app *Config) HandleGetPath(w http.ResponseWriter, r *http.Request) {
	path, err := app.GetPath()
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, path)
}

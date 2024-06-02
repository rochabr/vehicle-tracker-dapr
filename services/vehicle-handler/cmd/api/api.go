package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	Topic      = "routes"
	PubSubName = "vtd.pubsub"
)

// Handles the update loyalty endpoint
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

func (app *Config) HandleCreateShipment(w http.ResponseWriter, r *http.Request) {

	// Create order summary
	shipment, err := app.createShipment()
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// Send order to pub/sub
	ctx := context.Background()
	err = app.daprClient.PublishEvent(ctx, PubSubName, Topic, shipment)
	if err != nil {
		log.Printf("Error publishing the order summary: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Shipment %v created. Vehicle: %v. Driver: %v.", shipment.ShipmentID, shipment.Vehicle.VehicleID, shipment.Vehicle.Driver)

	app.writeJSON(w, http.StatusCreated, shipment.ShipmentID)
}

func (app *Config) createShipment() (Shipment, error) {

	// Load path from JSON file
	path, err := app.GetPath()
	if err != nil {
		return Shipment{}, err
	}

	// Initialize and return the order summary
	shipment := Shipment{
		ShipmentID: uuid.New().String(),
		Vehicle: Vehicle{
			VehicleID: 1,
			Driver:    "John Doe",
		},
		Path: path,
	}

	return shipment, nil
}

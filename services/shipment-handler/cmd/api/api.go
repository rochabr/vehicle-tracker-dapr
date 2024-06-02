package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	Topic              = "routes"
	PubSubName         = "vtd.pubsub"
	ShipmentStateStore = "vtd.shipment.state"
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

// Get shipment by Id handler
func (app *Config) HandleGetShipmentById(w http.ResponseWriter, r *http.Request) {
	shipmentId := chi.URLParam(r, "shipmentId")

	if shipmentId == "" {
		app.writeError(w, errors.New("Shipment ID not provided"), http.StatusNotFound)
		return
	}

	shipment, err := app.getShipment(shipmentId)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	if shipment.ShipmentID == "" {
		app.writeJSON(w, http.StatusNotFound, "Shipment not found")
		return
	}

	app.writeJSON(w, http.StatusOK, shipment)
}

// Creates shipment handler
func (app *Config) HandlePostShipment(w http.ResponseWriter, r *http.Request) {

	// Create shipment
	shipment, err := app.createShipment()
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// marshall shipment to save to state store
	data, err := json.Marshal(shipment)
	if err != nil {
		log.Printf("Error marshalling shipment. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// save shipment to state store
	err = app.daprClient.SaveState(context.Background(), ShipmentStateStore, shipment.ShipmentID, data, nil)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Shipment %v created. Vehicle: %v. Driver: %v.", shipment.ShipmentID, shipment.Vehicle.VehicleID, shipment.Vehicle.Driver)

	app.writeJSON(w, http.StatusCreated, shipment.ShipmentID)
}

// Delete shipment handler
func (app *Config) HandleDeleteShipment(w http.ResponseWriter, r *http.Request) {
	shipmentId := chi.URLParam(r, "shipmentId")

	if shipmentId == "" {
		app.writeError(w, errors.New("Shipment ID not provided"), http.StatusNotFound)
		return
	}

	// Get the current loyalty information for the customer from the loyalty state store
	err := app.daprClient.DeleteState(context.Background(), ShipmentStateStore, shipmentId, nil)
	if err != nil {
		app.writeError(w, errors.New("error deleting shipment from state store. Error: %v"), http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, "Shipment deleted")
}

// Start shipment handler
func (app *Config) HandleStartShipment(w http.ResponseWriter, r *http.Request) {
	shipmentId := chi.URLParam(r, "shipmentId")

	if shipmentId == "" {
		app.writeError(w, errors.New("Shipment ID not provided"), http.StatusNotFound)
		return
	}

	shipment, err := app.getShipment(shipmentId)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	if shipment.ShipmentID == "" {
		app.writeJSON(w, http.StatusNotFound, "Shipment not found")
		return
	}

	log.Printf("Starting shipment %v", shipment.ShipmentID)

	// Send path to pub/sub
	ctx := context.Background()

	if shipment.Path.Positions == nil || len(shipment.Path.Positions) == 0 {
		app.writeError(w, errors.New("Shipment path not found"), http.StatusNotFound)
		return
	}

	// Publish each point in the path to the pub/sub
	var shipmentPosition ShipmentPosition
	shipmentPosition.ShipmentID = shipment.ShipmentID

	for _, point := range shipment.Path.Positions {
		shipmentPosition.Position = point
		err = app.daprClient.PublishEvent(ctx, PubSubName, Topic, shipmentPosition)
		if err != nil {
			log.Printf("Error publishing the order summary: %v", err)
			app.writeError(w, err, http.StatusBadRequest)
			return
		}
		log.Printf("Vehicle position point: %v", point)
		time.Sleep(3 * time.Second)
	}

	log.Printf("Shipment %v arrived.", shipmentId)
	app.writeJSON(w, http.StatusOK, "Shipment arrived")
}

// Helper methods

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

func (app *Config) getShipment(shipmentId string) (Shipment, error) {

	// Get the current loyalty information for the customer from the loyalty state store
	item, err := app.daprClient.GetState(context.Background(), ShipmentStateStore, shipmentId, nil)
	if err != nil {
		//app.writeError(w, errors.New("error getting shipment from state store. Error: %v"), http.StatusNotFound)
		return Shipment{}, err
	}

	if item.Value == nil {
		//app.writeJSON(w, http.StatusNotFound, "Shipment not found")
		return Shipment{}, nil
	}

	var shipment Shipment
	err = json.Unmarshal(item.Value, &shipment)
	if err != nil {
		//app.writeError(w, err, http.StatusBadRequest)
		return Shipment{}, err
	}

	return shipment, nil
}

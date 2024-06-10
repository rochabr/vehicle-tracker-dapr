package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	ShipmentStateStore       = "vtd.shipment.state"
	PathHandlerServiceDaprId = "path-handler"
)

// Handles the health endpoint for Dapr
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5180"
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
func (app *Config) HandleCreateShipment(w http.ResponseWriter, r *http.Request) {

	// Create shipment
	shipment, err := app.createShipment()
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	//save shipment status
	err = app.saveShipmentStatus(shipment)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// // marshall shipment to save to state store
	// data, err := json.Marshal(shipment)
	// if err != nil {
	// 	log.Printf("Error marshalling shipment. Error: %v", err)
	// 	app.writeError(w, err, http.StatusBadRequest)
	// 	return
	// }

	// // save shipment to state store
	// err = app.daprClient.SaveState(context.Background(), ShipmentStateStore, shipment.ShipmentID, data, nil)
	// if err != nil {
	// 	app.writeError(w, err, http.StatusBadRequest)
	// 	return
	// }

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
	ctx := context.Background()

	shipmentId := chi.URLParam(r, "shipmentId")

	if shipmentId == "" {
		app.writeError(w, errors.New("Shipment ID not provided"), http.StatusNotFound)
		return
	}

	// Get shipment
	shipment, err := app.getShipment(shipmentId)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// Check if shipment is found
	if shipment.ShipmentID == "" {
		app.writeJSON(w, http.StatusNotFound, "Shipment not found")
		return
	}

	// Check if shipment path is empty
	if shipment.Path.Positions == nil || len(shipment.Path.Positions) == 0 {
		app.writeError(w, errors.New("shipment path not found"), http.StatusNotFound)
		return
	}

	shipment.Status = ShipmentStatusEnRoute

	// Save shipment status
	err = app.saveShipmentStatus(shipment)
	if err != nil {
		app.writeError(w, errors.New("error changing shipment status"), http.StatusBadRequest)
		return
	}

	log.Printf("Starting shipment %v", shipment.ShipmentID)

	// implement actor client stub
	vehicleActor := new(VehicleActorStub)
	vehicleActor.sId = shipmentId
	app.daprClient.ImplActorClientStub(vehicleActor)

	// Invoke user defined method GetUser with user defined param api.User and response
	// using default serializer type json
	_, err = vehicleActor.StartShipment(ctx, &shipment)
	if err != nil {
		app.writeError(w, errors.New("error starting shipment actor"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, shipment.ShipmentID)
}

// Helper methods
func (app *Config) createShipment() (Shipment, error) {

	// Gte calculated path
	path, err := app.getPath()
	if err != nil {
		return Shipment{}, err
	}

	// create shipment
	shipment := Shipment{
		ShipmentID: uuid.New().String(),
		Vehicle: Vehicle{
			VehicleID: 1,
			Driver:    "John Doe",
		},
		Path:   path,
		Status: ShipmentStatusPending,
	}

	return shipment, nil
}

func (app *Config) saveShipmentStatus(shipment Shipment) error {
	// marshall shipment to save to state store
	data, err := json.Marshal(shipment)
	if err != nil {
		log.Printf("Error marshalling shipment. Error: %v", err)
		return err
	}

	// save shipment to state store
	err = app.daprClient.SaveState(context.Background(), ShipmentStateStore, shipment.ShipmentID, data, nil)
	if err != nil {
		log.Printf("Error saving shipment. Error: %v", err)
		return err
	}

	return nil
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

func (app *Config) getPath() (Path, error) {

	ctx := context.Background()
	log.Printf("get path")
	response, err := app.daprClient.InvokeMethod(ctx, PathHandlerServiceDaprId, "path", "get")
	if err != nil {
		log.Printf("error calling service: %v", err)
		return Path{Positions: []Position{}},
			err
	}

	var path Path
	err = json.Unmarshal(response, &path)
	if err != nil {
		log.Printf("error unmarshalling response: %v", err)
		return Path{Positions: []Position{}},
			err
	}

	log.Print("Loaded path")
	return path, nil
}

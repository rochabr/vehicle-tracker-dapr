package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

const (
	Topic              = "locations"
	PubSubName         = "vtd.pubsub"
	LocationStateStore = "vtd.location.state"
	Route              = "locations"
)

// Handles the health endpoint for Dapr
func (app *Config) HandleHealthz(w http.ResponseWriter, r *http.Request) {

	//set app port
	daprHttpPort := "5380"
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

type JSONObj struct {
	PubsubName string `json:"pubsubName"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

type Result struct {
	Data any `json:"data"`
}

// Handles Dapr Endpoint and registers the subscription endpoint with the topic, pubsub and route
func (app *Config) HandleDaprEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Print("Received request from Dapr")
	jsonData := []JSONObj{
		{
			PubsubName: PubSubName,
			Topic:      Topic,
			Route:      Route,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Error Marshalling json data. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
	}

	log.Print("Writing response to Dapr")
	_, err = w.Write(jsonBytes)
	//err = app.writeJSON(w, http.StatusOK, "{}")
	if err != nil {
		log.Printf("Error writing json response. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
	}
}

// Handles the current position(location) endpoint
func (app *Config) HandleCurrentPosition(w http.ResponseWriter, r *http.Request) {

	// Unmarshall customer order
	var result Result
	err := app.readJSON(w, r, &result)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	s, err := json.Marshal(result.Data)
	if err != nil {
		log.Printf("Error marshalling shipment position. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	var shipmentPosition ShipmentPosition
	err = json.Unmarshal(s, &shipmentPosition)
	if err != nil {
		log.Printf("Error unmarshalling positon. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("Storing shipment position %v for shipment %v.", shipmentPosition.Position, shipmentPosition.ShipmentID)

	// marshall shipment to save to state store
	data, err := json.Marshal(shipmentPosition)
	if err != nil {
		log.Printf("Error marshalling shipment. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	// save shipment location to state store
	err = app.daprClient.SaveState(context.Background(), LocationStateStore, shipmentPosition.ShipmentID, data, nil)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, "Shipment position stored successfully")
}

// Handles the current position(location) endpoint
func (app *Config) HandleGetLastPosition(w http.ResponseWriter, r *http.Request) {

	shipmentId := chi.URLParam(r, "shipmentId")
	log.Printf("2 Getting last position for shipment %v.", shipmentId)

	if shipmentId == "" {
		app.writeError(w, errors.New("missing shipmentId on request"), http.StatusBadRequest)
	}

	log.Printf("Getting last position for shipment %v.", shipmentId)

	// save shipment location to state store
	item, err := app.daprClient.GetState(context.Background(), LocationStateStore, shipmentId, nil)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	if item.Value == nil {
		app.writeJSON(w, http.StatusNotFound, "Shipment Position not found")
		return
	}

	var shipmentPosition ShipmentPosition
	err = json.Unmarshal(item.Value, &shipmentPosition)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, shipmentPosition)
}

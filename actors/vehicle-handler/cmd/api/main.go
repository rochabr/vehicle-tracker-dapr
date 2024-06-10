package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"

	daprd "github.com/dapr/go-sdk/service/http"
)

const (
	PubSub        = "vtd.pubsub"
	LocationTopic = "locations"
)

func vehicleActorFactory() actor.ServerContext {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	return &VehicleActor{
		daprClient: client,
	}
}

type VehicleActor struct {
	actor.ServerImplBaseCtx
	daprClient dapr.Client
}

func (t *VehicleActor) Type() string {
	return "VehicleActor"
}

func main() {

	//set app port
	appPort := "7100"
	if value, ok := os.LookupEnv("APP_PORT"); ok {
		appPort = value
	}

	s := daprd.NewService(":" + appPort)
	s.RegisterActorImplFactoryContext(vehicleActorFactory)
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}

func StartShipment(shipment *Shipment) bool {
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if shipment == nil {
		log.Println("Shipment not found")
		return false
	}

	log.Println("Shipment path:", len(shipment.Path.Positions))

	for _, position := range shipment.Path.Positions {
		log.Println("Position:", position.Latitude, "-", position.Longitude)

		shipmentPosition := ShipmentPosition{
			ShipmentID: shipment.ShipmentID,
			Position:   position,
		}

		jsonString, err := json.Marshal(shipmentPosition)
		if err != nil {
			log.Printf("Error marshaling shipment position for Shipment: %s. Error: %v", shipment.ShipmentID, err)
			return false
		}
		log.Println("Object:", string(jsonString))

		err = client.PublishEvent(context.Background(), PubSub, LocationTopic, jsonString)
		if err != nil {
			log.Printf("Error publishing shipment position for Shipment: %s. Error: %v", shipment.ShipmentID, err)
			return false
		}
		log.Println("Published last position data:", shipmentPosition)

		time.Sleep(3 * time.Second)
	}

	return true
}

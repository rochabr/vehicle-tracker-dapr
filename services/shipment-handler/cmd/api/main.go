package main

import (
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
)

type Config struct {
	daprClient dapr.Client
}

func main() {
	log.Println("Starting the shipment handler.")

	//set app port
	appPort := "5100"
	if value, ok := os.LookupEnv("APP_PORT"); ok {
		appPort = value
	}

	// Initialize the Dapr client
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error creating dapr client: %v", err)
	}

	defer client.Close()

	app := Config{
		daprClient: client,
	}

	log.Printf("Starting the application on port %s\n", appPort)

	// create a new server
	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

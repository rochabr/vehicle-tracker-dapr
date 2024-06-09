package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	baseURL = "http://127.0.0.1:5100/shipment"
)

type Shipment struct {
	ID   string `json:"id"`
	To   string `json:"to"`
	From string `json:"from"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command: start, get, create, or delete")
		return
	}

	command := os.Args[1]
	switch command {
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a shipment ID")
			return
		}
		shipmentID := os.Args[2]
		startShipment(shipmentID)
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a shipment ID")
			return
		}
		shipmentID := os.Args[2]
		getShipment(shipmentID)
	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Please provide 'to' and 'from' addresses")
			return
		}
		to := os.Args[2]
		from := os.Args[3]
		createShipment(to, from)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a shipment ID")
			return
		}
		shipmentID := os.Args[2]
		deleteShipment(shipmentID)
	default:
		fmt.Println("Invalid command")
	}
}

func startShipment(shipmentID string) {
	url := fmt.Sprintf("%s/%s/start", baseURL, shipmentID)
	fmt.Printf("url: %v\n", url)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		fmt.Printf("Failed to start shipment: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to start shipment: %s\n", resp.Status)
		return
	}

	fmt.Printf("Shipment %s started successfully\n", shipmentID)
}

func getShipment(shipmentID string) {
	url := fmt.Sprintf("%s/%s", baseURL, shipmentID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to get shipment: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get shipment: %s\n", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	fmt.Println(string(body))
}

func createShipment(to, from string) {
	url := fmt.Sprintf("%s", baseURL)
	payload := fmt.Sprintf(`{"to": "%s", "from": "%s"}`, to, from)
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		fmt.Printf("Failed to create shipment: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Failed to create shipment: %s\n", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("Shipment created successfully with ID: %s\n", string(body))
}

func deleteShipment(shipmentID string) {
	url := fmt.Sprintf("%s/%s", baseURL, shipmentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("Failed to create delete request: %v\n", err)
		return
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to delete shipment: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to delete shipment: %s\n", resp.Status)
		return
	}

	fmt.Printf("Shipment %s deleted successfully\n", shipmentID)
}

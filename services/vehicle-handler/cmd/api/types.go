package main

import (
	"log"
	"os"
)

const FileName = "path.json"

type Shipment struct {
	ShipmentID string  `json:"shipmentId"`
	Vehicle    Vehicle `json:"vehicle"`
	Path       Path    `json:"route"`
}

type Vehicle struct {
	VehicleID int    `json:"vehicleId"`
	Driver    string `json:"driver"`
}

type Path struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Load path from JSON file
func (app *Config) GetPath() (Path, error) {

	fileName := ""
	if fn := os.Getenv("PATH_FILENAME"); fn != "" {
		fileName = fn
	} else {
		fileName = "cmd/api/paths/" + FileName
	}

	log.Printf("Filename %s\n", fileName)

	//fileName = "app/" + fileName

	var positions []Position

	path := Path{
		Positions: positions,
	}

	err := app.readJSONFile(fileName, &positions)
	if err != nil || len(positions) == 0 {
		log.Fatalf("Error loading  path from file %s. Returning empty path. Error: %s", fileName, err)
		return path, err
	}

	//load positions on path
	path.Positions = positions

	log.Printf("Loaded %d path from %s\n", len(positions), fileName)

	return path, nil
}

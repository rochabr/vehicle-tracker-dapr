package main

const (
	ShipmentStatusPending   = "pending"
	ShipmentStatusEnRoute   = "en-route"
	ShipmentStatusCompleted = "completed"
)

type Shipment struct {
	ShipmentID string  `json:"shipmentId"`
	Vehicle    Vehicle `json:"vehicle"`
	Path       Path    `json:"path"`
	Status     string  `json:"status"`
}

type Vehicle struct {
	VehicleID int    `json:"vehicleId"`
	Driver    string `json:"driver"`
}

type Path struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	Latitude  string `json:"_lat"`
	Longitude string `json:"_lon"`
}

type ShipmentPosition struct {
	ShipmentID string   `json:"shipmentId"`
	Position   Position `json:"position"`
}

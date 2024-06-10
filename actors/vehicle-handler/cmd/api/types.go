package main

import "context"

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
	Latitude  float64 `json:"_lat"`
	Longitude float64 `json:"_lon"`
}

type ShipmentPosition struct {
	ShipmentID string   `json:"shipmentId"`
	Position   Position `json:"position"`
}

type VehicleActorStub struct {
	sId           string
	StartShipment func(context.Context, *Shipment) (bool, error)
}

func (a *VehicleActorStub) Type() string {
	return "VehicleActor"
}

func (a *VehicleActorStub) ID() string {
	return a.sId
}

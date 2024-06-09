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
	Latitude  float64 `json:"_lat"`
	Longitude float64 `json:"_lon"`
}

type ShipmentPosition struct {
	ShipmentID string   `json:"shipmentId"`
	Position   Position `json:"position"`
}

type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
	Geometry   Geometry          `json:"geometry"`
}

type Geometry struct {
	Coordinates [][]float64 `json:"coordinates"`
	Type        string      `json:"type"`
}

// // ParseGeoJSON parses the given JSON string into a GeoJSON struct
// func ParseGeoJSON(jsonStr string) (*GeoJSON, error) {
// 	geoJSON := &GeoJSON{}
// 	err := json.Unmarshal([]byte(jsonStr), geoJSON)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return geoJSON, nil
// }

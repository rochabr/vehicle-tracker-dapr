package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON reads a JSON object from an HTTP request
func (app *Config) readJSONFile(jsonFile string, data any) error {
	file, err := os.Open(jsonFile)
	if err != nil {
		return err
	}

	// decode the JSON object
	dec := json.NewDecoder(file)
	err = dec.Decode(data)
	if err != nil {
		log.Fatalf("Error decoding JSON file: %v", err)
		return err
	}

	// ensure there are no additional bytes in the request body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		log.Fatalf("Error decoding JSON file: %v", err)
		return errors.New("body must only have a single JSON object")
	}

	return nil

}

// readJSON reads a JSON object from an HTTP request
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// decode the JSON object
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// ensure there are no additional bytes in the request body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON object")
	}

	return nil

}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) writeError(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

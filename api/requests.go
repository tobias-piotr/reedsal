package api

import (
	"encoding/json"
	"net/http"
)

// Decodes and processes a JSON payload from the HTTP request into passes payload pointer.
func ProcessPayload(w http.ResponseWriter, r *http.Request, payload Payload) error {
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(map[string]string{"error": "Invalid request payload"})
		w.Write(res)
		return err
	}
	return nil
}

// Runs validation method on passed payload and returns formatted error message in case of errors.
func ValidatePayload(w http.ResponseWriter, payload Payload) error {
	err := payload.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(ErrorResponse{"Invalid request payload", err.Details})
		w.Write(res)
		return err
	}
	return nil
}

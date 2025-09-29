package services

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Error string `json:"error"`
}

// writeJSONError writes a JSON error payload with the given HTTP status code.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{Error: message})
}

// writeJSON writes any payload as JSON with status 200 (or the current status if already set).
func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

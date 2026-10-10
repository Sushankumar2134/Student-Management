package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, status, ErrorResponse{Message: message})
}

func readJSON(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func getPathParam(r *http.Request, param string) string {
	return r.PathValue(param)
}
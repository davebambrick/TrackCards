package handlers

import (
	"encoding/json"
	"net/http"
)

// Health Render Handler
func Health(w http.ResponseWriter, r *http.Request) {
	object := Object{Status: "OK"}
	body, err := json.Marshal(object)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(body)
	if err != nil {

	}
}

type Object struct {
	Status string `json:"status"`
}

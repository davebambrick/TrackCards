package handlers

import (
	"encoding/json"
	"net/http"
)

// Default Render Handler
func Default(w http.ResponseWriter, r *http.Request) {
	object := DefaultObject{Value: "Hello World"}
	body, err := json.Marshal(object)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(body)
	if err != nil {

	}
}

type DefaultObject struct {
	Value string `json:"value"`
}

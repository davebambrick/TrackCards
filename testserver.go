package main

import (
	"net/http"

	"github.com/davebambrick/TrackCards/handlers"
)

func main() {
	http.HandleFunc("/", handlers.TrackHandler)
	http.ListenAndServe(":8900", nil)
}

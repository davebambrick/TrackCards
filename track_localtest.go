package main

import (
	"encoding/json"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	tf "github.com/davebambrick/TrackCards/format"
)

//// TEST OBJECT TYPES ////////////////////////////////
type Artist struct {
	CardToken  string      `json:"card_token"`
	ID         string      `json:"id"`
	ImageURL   string      `json:"image"`
	Name       string      `json:"name"`
	Popularity interface{} `json:"popularity"`
	Type       string      `json:"type"`
	Value      string      `json:"value"`
}

///////////////////////////////////////////////////////
type Album struct {
	Artist    *Artist `json:"artist"`
	CardToken string  `json:"card_token"`
	ID        string  `json:"id"`
	ImageURL  string  `json:"image"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Value     string  `json:"value"`
}

///////////////////////////////////////////////////////
type TrackEntity struct {
	Album          *Album      `json:"album"`
	AliasedFieldId string      `json:"aliased_field_id"`
	CardToken      string      `json:"card_token"`
	Duration       int         `json:"duration"`
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Popularity     interface{} `json:"popularity"`
	ResolveType    string      `json:"resolve_type"`
	Score          int         `json:"score"`
	SpokenName     string      `json:"spoken_name"`
	StreamURL      string      `json:"stream_url"`
	Type           string      `json:"type"`
	Value          string      `json:"value"`
}

///////////////////////////////////////////////////////

// Track Handler
func TrackHandler(w http.ResponseWriter, r *http.Request) {
	//// CREATE TEST ENTITY /////////////////////////////
	myArtist := Artist{
		"0000",
		"1111",
		"http://www.example.com/01",
		"Kendrick Lamar",
		0.9,
		"music_artist",
		"Kendrick Lamar",
	}
	///////////////////////////////////////////////////////
	myAlbum := Album{
		&myArtist,
		"2222",
		"3333",
		"http://www.example.com/02",
		"To Pimp a Butterfly",
		"music_album",
		"To Pimp a Butterfly",
	}
	///////////////////////////////////////////////////////
	myTrackEntity := TrackEntity{
		&myAlbum,
		"song",
		"4444",
		301,
		"5555",
		"Hood Politics",
		nil,
		"context",
		0,
		"Hood Politics by Kendrick Lamar",
		"http://www.example.com/03",
		"music_track",
		"Hood Politics by Kendrick Lamar",
	}
	///////////////////////////////////////////////////////

	// In the final version, the entity will be pulled from the request,
	// so we won't have to do this 2-step type conversion

	myTrackJson, _ := json.Marshal(myTrackEntity)       // Convert entity to json
	myTrackSimple, _ := simplejson.NewJson(myTrackJson) // Convert json to simplejson

	specList := tf.BuildSpecList(tf.TransformLibrary["Track"])  // Build spec list string
	transformed, _ := tf.TransformJSON(myTrackSimple, specList) // Transform w/ loaded Kazaam object
	writeable, _ := transformed.EncodePretty()                  // Encode as json for test printing

	w.Header().Set("Content-Type", "application/json")

	_, err := w.Write(writeable)
	if err != nil {
		log.Fatal("Could not write to the server")
	}
}

func main() {
	http.HandleFunc("/", TrackHandler)
	http.ListenAndServe(":8000", nil)
}

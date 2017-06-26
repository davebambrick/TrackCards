package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//// TEST OBJECT TYPES ////////////////////////////////
type ArtistEntity struct {
	CardToken  string      `json:"card_token"`
	ID         string      `json:"id"`
	ImageURL   string      `json:"image"`
	Name       string      `json:"name"`
	Popularity interface{} `json:"popularity"`
	Type       string      `json:"type"`
	Value      string      `json:"value"`
}

///////////////////////////////////////////////////////
type AlbumEntity struct {
	Artist    *ArtistEntity `json:"artist"`
	CardToken string        `json:"card_token"`
	ID        string        `json:"id"`
	ImageURL  string        `json:"image"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Value     string        `json:"value"`
}

///////////////////////////////////////////////////////
type TrackEntity struct {
	Album          *AlbumEntity `json:"album"`
	AliasedFieldId string       `json:"aliased_field_id"`
	CardToken      string       `json:"card_token"`
	Duration       int          `json:"duration"`
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Popularity     interface{}  `json:"popularity"`
	ResolveType    string       `json:"resolve_type"`
	Score          int          `json:"score"`
	SpokenName     string       `json:"spoken_name"`
	StreamURL      string       `json:"stream_url"`
	Type           string       `json:"type"`
	Value          string       `json:"value"`
}

func main() {
	//// CREATE TEST ENTITY /////////////////////////////
	myArtist := ArtistEntity{
		"0000",
		"1111",
		"http://www.example_url.com/01",
		"Kendrick Lamar",
		0.9,
		"music_artist",
		"Kendrick Lamar",
	}
	///////////////////////////////////////////////////////
	myAlbum := AlbumEntity{
		&myArtist,
		"2222",
		"3333",
		"http://www.example_url.com/02",
		"To Pimp a Butterfly",
		"music_album",
		"To Pimp a Butterfly",
	}
	///////////////////////////////////////////////////////
	myTrack := TrackEntity{
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
		"http://www.example_url.com/03",
		"music_track",
		"Hood Politics by Kendrick Lamar",
	}

	trackJson, _ := json.MarshalIndent(myTrack, "", " ")
	url := "http://localhost:8320"
	contentType := "application/json"
	resp, err := http.Post(url, contentType, bytes.NewReader(trackJson))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

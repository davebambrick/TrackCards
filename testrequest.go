package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Artist struct {
	CardToken  string      `json:"card_token"`
	ID         string      `json:"id"`
	ImageURL   string      `json:"image"`
	Name       string      `json:"name"`
	Popularity interface{} `json:"popularity"`
	Type       string      `json:"type"`
	Value      string      `json:"value"`
}

type Album struct {
	Artist    *Artist `json:"artist"`
	CardToken string  `json:"card_token"`
	ID        string  `json:"id"`
	ImageURL  string  `json:"image"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Value     string  `json:"value"`
}

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

func main() {
	Artist := Artist{
		"1234",
		"1234",
		"http://www.example.com/01",
		"Adele",
		0.8,
		"music_artist",
		"Adele",
	}
	Album := Album{
		&Artist,
		"2468",
		"2468",
		"http://www.example.com/02",
		"Hello",
		"music_album",
		"Hello",
	}
	Track := TrackEntity{
		&Album,
		"song",
		"1357",
		290,
		"1357",
		"Hello",
		nil,
		"context",
		0,
		"Hello by Adele",
		"http://www.example.com/03",
		"music_track",
		"Hello by Adele",
	}
	TrackJSON, _ := json.Marshal(Track)

	resp, _ := http.Post("http://localhost:8900", "application/json", bytes.NewReader(TrackJSON))
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Fatal(err)
	} else {
		var stuff map[string]interface{}
		json.Unmarshal(body, &stuff)
		indented, _ := json.MarshalIndent(stuff, "", " ")
		fmt.Println(string(indented))
	}
}

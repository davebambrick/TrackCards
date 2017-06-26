package main

import (
	"encoding/json"
	"log"
	"net/http"

	kz "gopkg.in/qntfy/kazaam.v3"
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

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	//// TEST ENTITY ////
	myArtist := Artist{
		"1234",
		"1234",
		"http://www.example.com/01",
		"Adele",
		0.8,
		"music_artist",
		"Adele",
	}
	myAlbum := Album{
		&myArtist,
		"2468",
		"2468",
		"http://www.example.com/02",
		"Hello",
		"music_album",
		"Hello",
	}
	myTrackEntity := TrackEntity{
		&myAlbum,
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
	TrackJson, err := json.Marshal(myTrackEntity)

	if err != nil {
		log.Fatal("Track object could not be converted to JSON format")
	}

	trackKazaam, err := kz.NewKazaam(`[{"operation":"shift", "spec": {"cardToken":"album.card_token","audioUrl":"stream_url","subtitle1":"album.artist.name","subtitle2":"album.name","title":"name","backgroundImageUrl":"album.image","extraData.trackInfo.durationInSeconds":"duration"}},{"operation":"default","spec": {"extraDataUrl":"http://www.example.com/04"}}]`)
	if err != nil {
		log.Fatal("Could not create new Kazaam instance")
	}

	JsonShift, err := trackKazaam.Transform(TrackJson)
	if err != nil {
		log.Fatal("Could not Transfrom Json")
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(JsonShift)
	if err != nil {
		log.Fatal("broken")
	}
}

func main() {
	http.HandleFunc("/", TrackHandler)
	http.ListenAndServe(":8000", nil)
}

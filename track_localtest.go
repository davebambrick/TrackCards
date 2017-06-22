package main

import (
	"encoding/json"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	kz "gopkg.in/qntfy/kazaam.v2"
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

// Track Handler
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
	trackSimple, err := simplejson.NewJson(TrackJson)
	if err != nil {
		log.Fatal("TrackJSON object could not be converted to simpleJSON format")
	}
	trackKazaam, err := kz.NewKazaam(`[{"operation":"shift", "spec": {"cardToken":"album.card_token","audioUrl":"stream_url","subtitle1":"album.artist.name","subtitle2":"album.name","title":"name","backgroundImageUrl":"album.image","extraData.trackInfo.durationInSeconds":"duration"}},{"operation":"default","spec": {"extraDataUrl":"http://www.example.com/04"}}]`)
	if err != nil {
		log.Fatal("Could not create new Kazaam instance")
	}

	transformedSimple, err := trackKazaam.Transform(trackSimple)
	if err != nil {
		log.Fatal("Could not transform simplejson")
	}
	transformedWriteable, err := transformedSimple.EncodePretty()
	if err != nil {
		log.Fatal("Could not marshal transformed simplejson")
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(transformedWriteable)
	if err != nil {
		log.Fatal("broken")
	}
}

func main() {
	http.HandleFunc("/", TrackHandler)
	http.ListenAndServe(":8000", nil)
}

///// TRANSFORMATIONS /////
var transforms = map[string]string{
	"cardToken":                             "album.card_token",
	"audioUrl":                              "stream_url",
	"subtitle1":                             "album.artist.name",
	"subtitle2":                             "album.name",
	"title":                                 "name",
	"backgroundImageUrl":                    "album.image",
	"extraData.trackInfo.durationInSeconds": " duration",
}

/*
"entity": {
          "album": {
            "artist": {
              "card_token": "142111",
              "id": "142111",
              "image": "http://artwork-cdn.7static.com/static/img/artistimages/00/001/421/0000142111_300.jpg",
              "name": "Adele",
              "popularity": 0.84,
              "type": "music_artist",
              "value": "Adele"
            },
            "card_token": "4872748",
            "id": "4872748",
            "image": "http://artwork-cdn.7static.com/static/img/sleeveart/00/048/727/0004872748_800.jpg",
            "name": "Hello",
            "type": "music_album",
            "value": "Hello"
          },
          "aliased_field_id": "song",
          "card_token": "49422899",
          "duration": 295,
          "id": "49422899",
          "name": "Hello",
          "popularity": null,
          "resolve_type": "context_named_entity_resolution",
          "score": 0,
          "spoken_name": "Hello by Adele from the album Hello",
          "stream_url": "http://iamplus-music-api-dev.herokuapp.com/stream_url?track_id=49422899",
          "type": "music_track",
          "value": "Hello by Adele from the album Hello"
}

{
    "cardToken": "60372376",
    "audioUrl": "http://iamplus-music-api-dev.herokuapp.com/stream_url?track_id=60372376",
    "subtitle1": "Power Music Workout",
    "subtitle2": "Workout Music Source - Kickbox Training Session (Non-Stop Workout Session 133-145 BPM)",
    "title": "Hello",
    "backgroundImageUrl": "http://artwork-cdn.7static.com/static/img/sleeveart/00/060/813/0006081338_800.jpg",
    "extraDataUrl": "http://54.190.11.200:8080/track/642227/60372376",
    "extraData": {
    	"trackInfo": {
    		"durationInSeconds": 219
    	}
    }
}


*/

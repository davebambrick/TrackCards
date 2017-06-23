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

///// TRANSFORMATIONS /////
var transforms = map[string]string{
	"cardToken":                             "entity.album.card_token",  // card ID
	"audioUrl":                              "entity.stream_url",        // audio streaming URL
	"subtitle1":                             "entity.album.artist.name", // artist name subtitle
	"subtitle2":                             "entity.album.name",        // album name subtitle
	"title":                                 "entity.name",              // track title
	"backgroundImageUrl":                    "entity.album.image",       // track image hosting URL
	"extraDataUrl":                          "http://www.music.com",     // metadata hosting URL
	"extraData.trackInfo.durationInSeconds": "entity.duration",          // duration of song
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

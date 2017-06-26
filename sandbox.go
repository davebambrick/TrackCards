package main

import (
	"encoding/json"
	"fmt"
  "strings"
)

type Specs map[string]interface{}

type Operation struct {
	Type  string
	Specs Specs
}
var TransformLibrary = map[string][]Operation{
	"music_track": []Operation{
		Operation{
			"shift",
			Specs{
				"cardToken":                             "album.card_token",
				"audioUrl":                              "stream_url",
				"subtitle1":                             "album.artist.name",
				"subtitle2":                             "album.name",
				"title":                                 "name",
				"backgroundImageUrl":                    "album.image",
				"extraData.trackInfo.durationInSeconds": "duration",
			},
		},
		Operation{
			"default",
			Specs{
				"extraDataUrl": "http://www.example_url.com/03",
			},
		},
	},
	"music_album":  []Operation{},
	"music_artist": []Operation{},
}

func main() {

	entityType, _ := json.Marshal("music_track")
  trimmed := strings.Trim(string(entityType),"\"")
  fmt.Println(trimmed)
	ops1 := TransformLibrary[trimmed]
	fmt.Println(ops1)


	stringType := "music_track"
  fmt.Println(stringType)
	ops2 := TransformLibrary[stringType]
	fmt.Println(ops2)

}

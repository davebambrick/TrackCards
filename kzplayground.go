package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	kz "gopkg.in/qntfy/kazaam.v2"
)

//// Content Object Types ///////////////
type Track struct {
	ArtistName string `json:"artist_name"`
	ImageURL   string `json:"image_url"`
	TrackName  string `json:"track_name"`
}

type TrackEntity struct {
	TrackInfo *Track `json:"track"`
	TrackID   string `json:"track_id"`
	StreamURL string `json:"streaming_url"`
}

//// Formatter Object Types //////////////
type Specs map[string]interface{}

type Operation struct {
	Type  string
	Specs Specs
}

func SpecParser(opType string, specs Specs) string {
	/*
		SpecParser takes the specifications of a given Kazaam operation,
		parses the specifications, and returns a string formatted as json.

		Args:
			opType (string): type of operation to eventually be performed by Kazaam
			spec (Specs): map[string]interface{} of operation-specific fields to their parameters

		Returns:
			json-formatted string of specifications according to operation-specific format
	*/
	var specString string
	switch opType {
	case "coalesce":
		for field, keyList := range specs {
			specString += "\"" + field + "\": [" + strings.Trim(strings.Join(keyList.([]string), " ,"), ",") + "]"
		}
	case "concat":
		//TODO
	case "extract":
		//TODO
	case "pass":
		//TODO
	default: //catches "shift" and "default" cases
		for target, source := range specs {
			specString += "\"" + target + "\": " + "\"" + source.(string) + "\","
		}
	}
	return "{" + strings.Trim(specString, ",") + "}"
}

func BuildSpecList(ops []Operation) string {
	/*
		BuildSpecList takes in a slice of Operations (type struct) and constructs a
		stringified list of operation specs formatted as json. Makes calls to
		SpecParser for each operation.

		Args:
			ops ([]Operation): slice of Operation objects (type struct) where each
			Operation contains its operation type (type string) and specifications
			(type map[string]interface{})

		Returns:
			stringified list of json-formatted operation specs
	*/
	var specList string

	for _, op := range ops {
		specList += "{\"operation\":\"" + op.Type + "\", \"spec\": " + SpecParser(op.Type, op.Specs) + "},"
	}
	return "[" + strings.Trim(specList, ",") + "]"
}

func TransformJSON(entity interface{}, specList string) (*simplejson.Json, error) {
	/*
		TransformJSON takes in an entity to be transformed, converts it to a
		simplejson.Json format, applies a Kazaam transformation to it in accordance
		with the provided specList, and returns a pointer to a new transformed
		simplejson.Json object.

		Args:
			entity (interface{}): entity to be converted to simplejson.Json and
			transformed via a Kazaam object
			specList (string): stringified list of json-formatted operation specifications,
			as supplied by BuildSpecList

		Returns:
			*simplejson.Json pointing to a transformation of the original entity input

	*/
	entityJson, err := json.Marshal(entity)
	if err != nil {
		return simplejson.New(), err
	}
	entitySimple, err := simplejson.NewJson(entityJson)
	if err != nil {
		return simplejson.New(), err
	}
	testKazaam, err := kz.NewKazaam(specList)
	if err != nil {
		return simplejson.New(), err
	}
	transformedSimple, err := testKazaam.Transform(entitySimple)
	if err != nil {
		return simplejson.New(), err
	}
	return transformedSimple, nil
}

func main() {

	Track := Track{
		"Adele",
		"http://artwork-cdn.7static.com/static/img/artistimages/00/001/421/0000142111_300.jpg",
		"Hello",
	}

	Entity := TrackEntity{
		&Track,
		"49422899",
		"http://iamplus-music-api-dev.herokuapp.com/stream_url?track_id=49422899",
	}
	//TODO: AUTOMATE THIS
	Op1 := Operation{
		"shift",
		Specs{
			"track_info":    "track",
			"id_of_track":   "track_id",
			"name_of_track": "track.track_name",
		},
	}
	Op2 := Operation{
		"default",
		Specs{
			"uniqueID": "1234567",
		},
	}
	// We can define the transformation of an entity with a list of operations
	ops := []Operation{Op1, Op2}
	specList := BuildSpecList(ops)

	transformedJSON, err := TransformJSON(Entity, specList)
	if err != nil {
		log.Fatal(err)
	}
	result, _ := transformedJSON.EncodePretty()
	fmt.Println(string(result))
}

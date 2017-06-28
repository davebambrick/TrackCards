// Package entityformatter provides an interface to
// create and utilize dynamic Kazaam operations to be
// performed on appropriate JSON objects.
package entityformatter

import (
	"errors"
	"fmt"
	"log"
	"strings"

	kz "gopkg.in/qntfy/kazaam.v3"
)

// Specs maps strings naming operation types to the given
// specifications for that operation.
type Specs map[string]interface{}

// Operation is a data structure representing a Kazaam operation
// in a more dynamic format. In this structure we can edit operations
// with ease.
type Operation struct {
	Type  string
	Specs Specs
}

// TransformLibrary is the current structure we are using
// to store the set of operations for a given entity. Entity types
// as strings are mapped to a sequence (slice) of Operation objects
// representing the complete transformation to be performed by Kazaam.
var TransformLibrary = map[string][]Operation{
	"music_track": []Operation{
		Operation{
			"concat",
			Specs{
				"sources": []map[string]string{
					{"value": "http:/"},
					{"value": "host_goes_here"},
					{"value": "track"},
					{"path": "album.artist.id"},
					{"path": "id"},
				},
				"targetPath": "extraDataUrl",
				"delim":      "/",
			},
		},
		Operation{
			"shift",
			Specs{
				"cardToken":                             "id",
				"audioUrl":                              "stream_url",
				"subtitle1":                             "album.artist.name",
				"subtitle2":                             "album.name",
				"title":                                 "name",
				"backgroundImageUrl":                    "album.image",
				"extraData.trackInfo.durationInSeconds": "duration",
				"extraDataUrl":                          "extraDataUrl",
			},
		},
	},
	"\"music_album\"":  []Operation{},
	"\"music_artist\"": []Operation{},
}

// SpecParser takes the specifications of a given Operation Object,
// parses the specifications, and returns a string formatted as json.
//
// Args:
// 	opType (string): type of operation to eventually be performed by Kazaam
// 	spec (Specs): map[string]interface{} of operation-specific fields to their parameters
//
// Returns:
// 	json-formatted string of specifications according to operation-specific format
func SpecParser(op Operation) (string, error) {

	var specString string
	switch op.Type {

	case "coalesce":
		for specField, keyList := range op.Specs {
			specString += fmt.Sprintf("\"%s\": [%s], ", specField, strings.Trim(strings.Join(keyList.([]string), ", "), ", "))
		}

	case "concat":
		for specField, specVal := range op.Specs {
			switch specField {
			case "sources": // specValue is a list of values to be concatenated
				var concatString string
				for _, concat := range specVal.([]map[string]string) {
					for k, v := range concat {
						concatString += fmt.Sprintf("{\"%s\":\"%s\"}, ", k, v)
					}
				}
				specString += fmt.Sprintf("\"%s\": [%s], ", specField, strings.Trim(concatString, ", "))
			default:
				specString += fmt.Sprintf("\"%s\": \"%s\", ", specField, specVal)
			}
		}

	case "extract":
		specString += fmt.Sprintf("\"path\": \"%s\"", op.Specs["path"])

	case "pass":
		break
	default: //catches "shift" and "default" cases
		for target, source := range op.Specs {
			specString += fmt.Sprintf("\"%s\": \"%s\", ", target, source)
		}
	}
	specs := strings.Trim(specString, ", ")
	if op.Type != "pass" && len(specs) == 0 {
		return "", errors.New("Could not parse specifications")
	}
	return fmt.Sprintf("{\"operation\": \"%s\", \"spec\": {%s}}", op.Type, specs), nil
}

// BuildSpecList takes in a string stringType, pulls the corresponding sequence
// of operations from the TransformLibrary, and constructs a stringified list of
// operation specs formatted as json. Makes calls to SpecParser for each operation.
//
// Args:
// 	stringType (string): string of entity type to be transformed; we use this
// 	string as a key in the TransformLibrary to access the desired sequence of
// 	operations and specifications
//
// Returns:
// 	stringified list of json-formatted operation specs
func BuildSpecList(stringType string) (string, error) {

	var specList string
	for _, op := range TransformLibrary[stringType] {
		parsed, err := SpecParser(op)
		if err != nil {
			log.Fatal(err)
		}
		specList += parsed + ", "
	}
	return fmt.Sprintf("[%s]", strings.Trim(specList, ", ")), nil
}

// TransformJSON takes in a json entity to be transformed, applies a Kazaam
// transformation to it in accordance with the provided specList string, and
// returns the new transformed json object.
//
// Args:
// 	entity (interface{}): entity to be converted to simplejson.Json and
// 	transformed via a Kazaam object
// 	specList (string): stringified list of json-formatted operation specifications,
// 	as supplied by BuildSpecList
//
// Returns:
// 	transformation of the original entity json object as a slice of bytes
func TransformJSON(entity []byte, specList string) ([]byte, error) {

	kazaam, err := kz.NewKazaam(specList)
	if err != nil {
		return []byte{}, err
	}
	transformed, err := kazaam.Transform(entity)
	if err != nil {
		return []byte{}, err
	}
	return transformed, nil
}

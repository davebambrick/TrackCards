package format

import (
	"errors"
	"fmt"
	"log"
	"strings"

	kz "gopkg.in/qntfy/kazaam.v3"
)

//// Formatter Object Types //////////////
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

func SpecParser(op Operation) (string, error) {
	/*
		SpecParser takes the specifications of a given Operation Object,
		parses the specifications, and returns a string formatted as json.

		Args:
			opType (string): type of operation to eventually be performed by Kazaam
			spec (Specs): map[string]interface{} of operation-specific fields to their parameters

		Returns:
			json-formatted string of specifications according to operation-specific format
	*/
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
				for location, concat := range specVal.(map[string]string) {
					concatString += fmt.Sprintf("{\"%s\": \"%s\"}, ", location, concat)
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

func BuildSpecList(stringType string) (string, error) {
	/*
		BuildSpecList takes in a string stringType, pulls the corresponding sequence
		of operations from the TransformLibrary, and constructs a stringified list of
		operation specs formatted as json. Makes calls to SpecParser for each operation.

		Args:
			stringType (string): string of entity type to be transformed; we use this
			string as a key in the TransformLibrary to access the desired sequence of
			operations and specifications

		Returns:
			stringified list of json-formatted operation specs
	*/

	var specList string
	ops := TransformLibrary[stringType]
	for _, op := range ops {
		parsed, err := SpecParser(op)
		if err != nil {
			log.Fatal(err)
		}
		specList += parsed + ", "
	}
	return fmt.Sprintf("[%s]", strings.Trim(specList, ", ")), nil
}

func TransformJSON(entity []byte, specList string) ([]byte, error) {
	/*
		TransformJSON takes in a json entity to be transformed, applies a Kazaam
		transformation to it in accordance with the provided specList string, and
		returns the new transformed json object.

		Args:
			entity (interface{}): entity to be converted to simplejson.Json and
			transformed via a Kazaam object
			specList (string): stringified list of json-formatted operation specifications,
			as supplied by BuildSpecList

		Returns:
			transformation of the original entity json object as a slice of bytes

	*/

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

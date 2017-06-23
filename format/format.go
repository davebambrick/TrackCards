package reformat

import (
	"errors"
	"fmt"
	"log"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	kz "gopkg.in/qntfy/kazaam.v2"
)

//// Formatter Object Types //////////////
type Specs map[string]interface{}

type Operation struct {
	Type  string
	Specs Specs
}

var TransformLibrary = map[string][]Operation{
	"Track": []Operation{
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
				"extraDataUrl": "http://www.example_url.com/",
			},
		},
	},
	"Album":  []Operation{},
	"Artist": []Operation{},
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
		parsed, err := SpecParser(op)
		if err != nil {
			log.Fatal(err)
		}
		specList += parsed + ", "
	}
	return fmt.Sprintf("[%s]", strings.Trim(specList, ", "))
}

func TransformJSON(entity *simplejson.Json, specList string) (*simplejson.Json, error) {
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
			*simplejson.Json pointing to a new transformation of the original entity input

	*/

	testKazaam, err := kz.NewKazaam(specList)
	if err != nil {
		return simplejson.New(), err
	}

	transformedSimple, err := testKazaam.Transform(entity)
	if err != nil {
		return simplejson.New(), err
	}

	return transformedSimple, nil
}

package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	sj "github.com/bitly/go-simplejson"
	ef "github.com/davebambrick/TrackCards/entityformatter"
)

// TrackHandler takes in an HTTP POST request and pulls a json music entity object
// from the request body. It pulls the entity's type, and then performs a sequence
// of corresponding Kazaam transformations on the json, before writing the output
// transformed json to the ResponseWriter.
//
// Args:
// 	w (http.ResponseWriter): ResponseWriter object used to build the response
// 	to be sent to the client
// 	r (*http.Request): Pointer to the incoming Request object; it is intended
// 	to be a POST request so we can access the request body
//
// Returns:
// 	None
func TrackHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	trackJSON, err := ioutil.ReadAll(r.Body) // Read target entity from request body
	if err != nil {
		log.Fatal(err)
	}
	trackSimple, err := sj.NewJson(trackJSON) // Transform to simplejson to access entity type
	if err != nil {
		log.Fatal(err)
	}
	entityType, err := trackSimple.Get("type").Encode() // Pull entity type from simplejson
	if err != nil {
		log.Fatal(err)
	}
	specList, err := ef.BuildSpecList(strings.Trim(string(entityType), "\"")) // Pass entity type into spec builder, build spec list string
	if err != nil {
		log.Fatal(err)
	}
	transformed, err := ef.TransformJSON(trackJSON, specList) // Transform w/ loaded Kazaam object
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(transformed)
	if err != nil {
		log.Fatal("Could not write to the server")
	}
}

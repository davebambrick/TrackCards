package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	sj "github.com/bitly/go-simplejson"
	ef "github.com/davebambrick/TrackCards/entityformatter"
)

func TrackHandler(w http.ResponseWriter, r *http.Request) {

	///////////////////////////////////////////////////////
	defer r.Body.Close()

	trackJson, _ := ioutil.ReadAll(r.Body) // Read target entity from request body

	trackSimple, _ := sj.NewJson(trackJson)               // Transform to simplejson to access entity type
	entityType, _ := trackSimple.Get("type").Encode()     // Pull entity type from simplejson
	stringType := strings.Trim(string(entityType), " \"") // Pass entity type into spec builder

	specList, _ := ef.BuildSpecList(stringType)             // Build spec list string
	transformed, _ := ef.TransformJSON(trackJson, specList) // Transform w/ loaded Kazaam object

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(transformed)
	if err != nil {
		log.Fatal("Could not write to the server")
	}
}

package main

import(
	"fmt"
	"encoding/json"
	kz "gopkg.in/qntfy/kazaam.v2"
	simplejson "github.com/bitly/go-simplejson"
)

type Track struct {
	ArtistName string `json:"artist_name"`
	ImageURL string `json:"image_url"`
	TrackName string `json:"track_name"`
}
	
type TrackEntity struct {
	TrackInfo Track `json:"track"`
	TrackID string `json:"track_id"`
	StreamURL string `json:"streaming_url"`
	
}

func main() {
	
	myTrack := Track{
		"Adele",
		"http://artwork-cdn.7static.com/static/img/artistimages/00/001/421/0000142111_300.jpg",
		"Hello",
	}
	
	myEntity := TrackEntity{
		myTrack,
		"49422899",
		"http://iamplus-music-api-dev.herokuapp.com/stream_url?track_id=49422899",
	}

	fmt.Println("#################")
	fmt.Println("### RAWENTITY ###")
	fmt.Println("#################")
	fmt.Println(myEntity)
	fmt.Println("##################")
	fmt.Println("### SIMPLEJSON ###")
	fmt.Println("##################")
	myJson, _ := json.Marshal(myEntity)
	mySimple, _ := simplejson.NewJson(myJson)
	fmt.Println(*mySimple)
	fmt.Println("##################")
	fmt.Println("### KAZAAM'D 1 ###")
	fmt.Println("##################")
	myKazaam, _ := kz.NewKazaam(`[{"operation": "shift", "spec": {"id_of_track": "track_id", "url":"streaming_url", "name_of_track":"track.track_name","track_map":"track"}}]`)
	transformedJs, _ := myKazaam.Transform(mySimple)
	fmt.Println(*transformedJs)
	fmt.Println("##################")
	fmt.Println("### KAZAAM'D 2 ###")
	fmt.Println("##################")
	concatRawstr := `[{"operation": "concat","spec": {"sources": [{"value": "The song is called"},{"path": "name_of_track"},{"value": "by artist"},{"path":"track_map.artist_name"}],"targetPath": "name_of_track", "delim": "//"}}]`
	myKazaam2, _ := kz.NewKazaam(concatRawstr)
	transformedJs2, _ := myKazaam2.Transform(transformedJs)
	fmt.Println(*transformedJs2)
	fmt.Println("#######################")
	fmt.Println("### DOUBLE KAZAAM'D ###")
	fmt.Println("#######################")
	myKazaam3, _ := kz.NewKazaam(`[{"operation": "shift", "spec": {"id_of_track": "track_id", "url":"streaming_url", "name_of_track":"track.track_name","track_map":"track"}}, {"operation": "concat","spec": {"sources": [{"value": "The song is called"},{"path": "name_of_track"},{"value": "by artist"},{"path":"track_map.artist_name"}],"targetPath": "name_of_track", "delim": "//"}}]`)
	doubleTransformed, _ := myKazaam3.Transform(mySimple)
	fmt.Println(*doubleTransformed)
	
	/*
	newJs := json.
	newSimpleJs := simplejson.New()
	myKz, _ := kz.NewKazaam("")
	fmt.Println(*newJs)
	fmt.Println(*myKz.Transform(newJs))
	*/
}



/*
concat format
{
    "operation": "concat",
    "spec": {
        "sources": [{
            "value": "TEST"
        }, {
            "path": "a.timestamp"
        }],
        "targetPath": "a.timestamp",
        "delim": ","
    }
}
*/















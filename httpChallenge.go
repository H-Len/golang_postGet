/*- Write an HTTP server that functions as a key/value database.

Users will make a POST request to /db to create or assign a value.
User will make a GET request /db/<key> to get the value.
Both request and response will be in JSON format.

Use a map from string to the empty interface to hold values.
And since we can't access the Go data structure from two different Go routines, limit the access to the map with a sync.Mutex.

You can see an example of a POST body below. The key is x, and the value is 1.*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Data is the key and matching value from Json
type Data struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func main() {
	data := &Data{}
	// func for handler that does operation (dbHandler)
	dbPost := func(w http.ResponseWriter, r *http.Request) {
		// Decode request
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		req := &Data{}

		if err := dec.Decode(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data.Key = req.Key
		data.Value = req.Value
		fmt.Println(req)
	}

	dbGet := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(w)
		if err := enc.Encode(data); err != nil {
			// Can't return error to client here
			log.Printf("can't encode %v - %s", data, err)
		}

	}

	data.Key = "frangipani"
	data.Value = 7

	http.HandleFunc("/db2", dbGet)
	http.HandleFunc("/db", dbPost)
	if err := http.ListenAndServe("127.0.0.1:8088", nil); err != nil {
		log.Fatal(err)
	}

}

/*- Write an HTTP server that functions as a key/value database.

And since we can't access the Go data structure from two different Go routines, limit the access to the map with a sync.Mutex.

You can see an example of a POST body below. The key is x, and the value is 1.*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Data is the key and matching value from Json
type Data struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func main() {
	var db = make(map[string]int)

	dbPost := func(w http.ResponseWriter, r *http.Request) {
		// Decode request
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		req := &Data{}

		if err := dec.Decode(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db[req.Key] = req.Value

		fmt.Println(req)
	}

	dbGet := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer r.Body.Close()

		path := r.URL.EscapedPath()

		enc := json.NewEncoder(w)

		pathURL := strings.Split(path, "/db/")
		keySlice := pathURL[0:]                 //obtain key
		resultKey := strings.Join(keySlice, "") //convert array to string
		fmt.Println(resultKey)
		// look at map documentation and see if there's a method to check if key exists
		if val, ok := db[resultKey]; ok {
			fmt.Println("value: ", val)

			fmt.Println(db[resultKey])

			resp := &Data{resultKey, db[resultKey]}
			if err := enc.Encode(resp); err != nil {
				// Can't return error to client here
				log.Printf("can't encode %v - %s", resp, err)
			}
			fmt.Println(resp)
		} else {
			fmt.Println("key not found")
		}

	}

	// func for handler that does operation (dbHandler)
	dbHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dbGet(w, r)
		case http.MethodPost:
			dbPost(w, r)
		default:
			println("method not supported")
		}
	}

	// db["lasagna"] = 78
	// fmt.Println(db["lasagna"])

	http.HandleFunc("/db", dbHandler)
	http.HandleFunc("/db/", dbGet)

	// path := "/db/confusion"
	// fmt.Println(path)

	if err := http.ListenAndServe("127.0.0.1:8088", nil); err != nil {
		log.Fatal(err)
	}

}

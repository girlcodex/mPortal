package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"reflect"
	"strings"
)

func searchESAPI() {

	// !!todo!! bring the search terms in from an input
	// static test inputs in var declaration in main.go are:
	//           fieldName = "title"
	//           searchTerm = "scrum"

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	// Exit the system if connection raises an error
	if err != nil {
		log.Fatalf("Elasticsearch connection error:", err)
	}

	var mapResp map[string]interface{}
	var buf bytes.Buffer

	// Concatenate a string from query for reading
	var b strings.Builder
	b.WriteString(query)
	read := strings.NewReader(b.String())

	fmt.Println("read:", read)
	fmt.Println("read TYPE:", reflect.TypeOf(read))
	fmt.Println("JSON encoding:", json.NewEncoder(&buf).Encode(read))

	// Attempt to encode the JSON query and look for errors
	if err := json.NewEncoder(&buf).Encode(read); err != nil {
		log.Fatalf("Error encoding query: %s", err)

		// Query is a valid JSON object
	} else {
		fmt.Println("json.NewEncoder encoded query:", read, "\n")

		// Pass the JSON query to the Golang client's Search() method
		res, err := client.Search(
			client.Search.WithContext(ctx),
			client.Search.WithIndex(searchTerm),
			client.Search.WithBody(read),
			client.Search.WithTrackTotalHits(true),
			client.Search.WithPretty(),
		)

		// Check for any errors returned by API call to Elasticsearch
		if err != nil {
			log.Fatalf("Elasticsearch Search() API ERROR:", err)
			// If no errors are returned, parse esapi.Response object
		} else {
			fmt.Println("res TYPE:", reflect.TypeOf(res))

			// Close the result body when the function call is complete
			defer res.Body.Close()

			// Decode the JSON response and set a pointer
			if err := json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)

				// If no error, then convert response to a map[string]interface
			} else {
				//!!todo!
				fmt.Println("ERROR HERE:", mapResp, "\n")

				fmt.Println("mapResp TYPE:", reflect.TypeOf(mapResp), "\n")

				// Iterate the document "hits" returned by API call
				for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {

					// Parse the attributes/fields of the document
					doc := hit.(map[string]interface{})
					// The "_source" data is another map interface nested inside of doc
					source := doc["_source"]
					fmt.Println("doc _source:", reflect.TypeOf(source))

					// Get the document's _id and print it out along with _source data
					docID := doc["_id"]
					fmt.Println("docID:", docID)
					fmt.Println("_source:", source, "\n")
				} // end of response iteration

			}
		}
	}
}

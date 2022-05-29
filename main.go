package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/elastic/go-elasticsearch/v7"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Employee struct {
	uniqid string `json:"UUID"`
	empid  string `json:"Employee ID"`
	fname  string `json:"First Name"`
	lname  string `json:"Last Name"`
}

type Task struct {
	unitid     string `json:"GUID"`
	taskid     string `json:"Task ID"`
	assignedto string `json:"Assigned To"`
	title      string `json:"Title"`
	privacy    string `json:"Privacy"`
}

type Whois struct {
	unitid     string `json:"GUID"`
	taskid     string `json:"Task ID"`
	assignedto string `json:"Assigned To"`
	title      string `json:"Title"`
	privacy    string `json:"Privacy"`

	uniqid string `json:"UUID"`
	empid  string `json:"Employee ID"`
	fname  string `json:"First Name"`
	lname  string `json:"Last Name"`
}

type Message struct {
	Table  string           `json:"table"`
	Id     int              `json:"id"`
	Action string           `json:"action"`
	Data   *json.RawMessage `json:"data"`
}

// do better authentication than this, obviously
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "workers"
)

var (
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//
	elasticIndex = "https://cointest.ent.us-central1.gcp.cloud.es.io:9243"
	//
	ctx = context.Background()
	//
	cfg = elasticsearch.Config{
		CloudID: "cointest:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvJDcxY2I3N2ZkYTUwZDQxMjQ4YWQzNDlkMTFjMzAyZmEzJGYyNTQwMmJiNTdkMzQ0NTg5NmY2NmM0MTY3MDE2YTJj",
		APIKey:  "emR5ekI0RUJKUlZnRGc2QURHX286bFVYSnBoZXVUd3Fqel9pek5aQnplQQ==",
	}
	//
	fieldName  = "title" // static test input for searchESAPI.go
	searchTerm = "scrum" // static test input
	//
	query = `{"query": {"term": {"` + fieldName + `" : "` + searchTerm + `"}}, "size": 10}`
	//
	verbose = true
	//		inserts, deletes int64
	//		idRef            string

)

func main() {

	// ES client stuff

	searchESAPI()

	//es, _ := elasticsearch.NewClient(cfg)
	//check that it works
	//log.Println(elasticsearch.Version)

	//res, err := es.Info()
	//if err != nil {
	//	log.Fatalf("Error getting response: %s", err)
	//}

	//log.Println(res)

	//elasticListener()
	//	!TODO!

	// these are the api endpoints
	// endpoints are
	//
	//     search employees by uuid                                            //        /employees?uuid='b98291a1-69e9-4030-9afd-fd23a4d93f0f'
	//      ^^ this is relevent, i get it and its there, but in reality
	//         we'd use empid for humans to search, so:				           // 		 /employees?empid='1'
	//
	//	   search task by name								                   //        /tasks?name='scrum'
	//
	//	   search all employees working on a specific task			           //        /tasks?name=
	//      ^^ involves referencing the employees table empid

	http.HandleFunc("/employees", EMP)
	http.HandleFunc("/tasks", TAS)
	http.HandleFunc("/whois", WHO)
	http.HandleFunc("/newEmployee", empPOST)
	//	http.HandleFunc("/newTask", POSTHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

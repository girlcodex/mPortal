package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
)

type Employee struct {
	Uniqid string //`json:"UUID"`
	Empid  string //`json:"Employee ID"`
	Fname  string //`json:"First Name"`
	Lname  string //`json:"Last Name"`
}

type Task struct {
	Unitid     string //`json:"GUID"`
	Taskid     string //`json:"Task ID"`
	Assignedto string //`json:"Assigned To"`
	Title      string //`json:"Title"`
	Privacy    string //`json:"Privacy"`
}

type Whois struct {
	Uniqid     string //`json:"Task UUID"`
	Empid      string //`json:"Employee UUID"`
	Fname      string //`json:"Task ID"`
	Lname      string //`json:"Assigned To"`
	Unitid     string //`json:"Title"`
	Taskid     string //`json:"Privacy"`
	Assignedto string //`json:"Employee ID"`
	Title      string //`json:"First Name"`
	Privacy    string //`json:"Last Name"`
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
		CloudID: "",
		APIKey:  "",
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
	employeeData = []byte(`
		[
		  {
			"uniqid": "b98291a1-69e9-4030-9afd-fd23a4d93f0f",
			"empid": 1,
			"fname": "jane",
			"lname": "smith"
		  },
		  {
			"uniqid": "22f2ece1-bccc-47eb-b28c-9bc767f3dc89",
			"empid": 2,
			"fname": "billy",
			"lname": "jones"
		  },
		  {
			"uniqid": "87fab4fb-441f-4cff-9121-2e3222ea7d18",
			"empid": 3,
			"fname": "lee",
			"lname": "irving"
		  },
		  {
			"uniqid": "f056f333-7836-471f-8578-47b24fd8b911",
			"empid": 4,
			"fname": "sarah",
			"lname": "pilsner"
		  },
		  {
			"uniqid": "72064bd1-b5b5-4911-b177-42609de697f9",
			"empid": 5,
			"fname": "guy",
			"lname": "young"
		  },
		  {
			"uniqid": "2abe1d24-72ca-4f03-9f7b-7478d025f716",
			"empid": 6,
			"fname": "lady",
			"lname": "oldman"
		  }
		]`)

	tasksData = []byte(`
		[
		  {
			"unitid": "18f16020-0e13-4f16-b4a4-e8dd52237d51",
			"taskid": 1,
			"assignedto": [1, 2, 3, 4, 5, 6],
			"title": "scrum meeting",
			"privacy": 0
		  },
		  {
			"unitid": "ca8af363-cbd5-4416-87f1-d65da242f69d",
			"taskid": 2,
			"assignedto": [4, 5, 6],
			"title": "interview",
			"privacy": 0
		  },
		  {
			"unitid": "2e819dbc-c15d-4e4f-896f-c0a95bd19b42",
			"taskid": 3,
			"assignedto": [1, 2, 6],
			"title": "documentation",
			"privacy": 0
		  },
		  {
			"unitid": "a9b89431-d8d6-41d2-be87-76c00ffe8484",
			"taskid": 4,
			"assignedto": [1, 3],
			"title": "secret docker file generation",
			"privacy": 1
		  },
		  {
			"unitid": "0a284be9-4033-4c53-8d6f-2b50e2ee5342",
			"taskid": 5,
			"assignedto": [2],
			"title": "a/b testing",
			"privacy": 0
		  },
		  {
			"unitid": "082b65f8-336a-4aa8-9869-6b1885441e4f",
			"taskid": 6,
			"assignedto": [1, 2],
			"title": "secret scrum meeting",
			"privacy": 1
		  },
		  {
			"unitid": "e5c2b757-f4c3-435a-b387-a231b158466a",
			"taskid": 7,
			"assignedto": [6],
			"title": "push dev to prod",
			"privacy": 0
		  },
		  {
			"unitid": "551fac4c-088e-402e-b597-4c1ce7f67cfe",
			"taskid": 8,
			"assignedto": [4, 6],
			"title": "secret scrum meeting",
			"privacy": 1
		  },
		  {
			"unitid": "82ded647-75bf-4a39-b071-fd80c4c95fd8",
			"taskid": 9,
			"assignedto": [],
			"title": "vacation",
			"privacy": 0
		  }
		]
		`)
)

func main() {
	mainPSQL()
	// ES client stuff

	//searchESAPI()

	es, _ := elasticsearch.NewClient(cfg)
	//check that it works
	log.Println(elasticsearch.Version)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)
	handleRequests()
	elasticListener()
	//	!TODO!

}

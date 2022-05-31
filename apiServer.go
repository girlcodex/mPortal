package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "io/ioutil"
	"log"
	"net/http"
)

var employee []Employee

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Please specify an api endpoint")
	fmt.Println("{host}:/")
}

func returnAllEmployees(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/employees was called! apiServer.go:19 returnAllEmployees")
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM employees")
	if err != nil {
		log.Fatal(err)
	}

	var employees []Employee

	for rows.Next() {
		var employee Employee
		rows.Scan(&employee.Uniqid, &employee.Empid, &employee.Fname, &employee.Lname)
		employees = append(employees, employee)
	}
	empBytes, _ := json.MarshalIndent(employees, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(empBytes)

	defer rows.Close()
	defer db.Close()

}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["empid"]
	fmt.Println("/employees was called! apiServer.go:48 returnSingleEmployee     key=: " + key)
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM employees WHERE empid = " + key)
	if err != nil {
		log.Fatal(err)
	}

	var employees []Employee

	for rows.Next() {
		var employee Employee
		rows.Scan(&employee.Uniqid, &employee.Empid, &employee.Fname, &employee.Lname)
		employees = append(employees, employee)
	}
	empBytes, _ := json.MarshalIndent(employees, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(empBytes)

	defer rows.Close()
	defer db.Close()

}
func returnSingleByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := string(vars["uniqid"])
	fmt.Println("/employees was called! apiServer.go:75 returnSingleByUUID     key=: " + key)
	db := OpenConnection()
	//var temp1 = stri
	rows := db.QueryRow("SELECT * FROM employees WHERE uniqid::text = $1", key)
	//fmt.Println(rows)
	//if err != nil {
	//	log.Fatal(err)
	//}
	var employees []Employee

	var employee Employee
	rows.Scan(&employee.Uniqid, &employee.Empid, &employee.Fname, &employee.Lname)
	employees = append(employees, employee)

	empBytes, _ := json.MarshalIndent(employees, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(empBytes)

	//defer rows.Close()
	defer db.Close()

}
func returnSingleTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := string("%%" + vars["title"] + "%%")
	fmt.Println("/employees was called! apiServer.go:101 returnSingleByTitle     key=: " + key)
	db := OpenConnection()
	//var temp1 = stri
	rows, err := db.Query("SELECT * FROM tasks WHERE title like $1 and privacy = 0", key)
	//fmt.Println(rows)
	if err != nil {
		log.Fatal(err)
	}
	var tasks []Task
	for rows.Next() {
		var task Task
		rows.Scan(&task.Unitid, &task.Taskid, &task.Assignedto, &task.Title, &task.Privacy)
		tasks = append(tasks, task)
	}
	empBytes, _ := json.MarshalIndent(tasks, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(empBytes)

	//defer rows.Close()
	defer db.Close()

}

// using gorilla mux to normalize endpoint urls
func handleRequests() {
	// these are the api endpoints
	// endpoints are
	//
	//     search employees by uuid                             //       /employees?uuid='b98291a1-69e9-4030-9afd-fd23a4d93f0f'
	//     search employees by empid				         	// 		 /employees?empid='1'
	//	   search task by name								  	//       /tasks?name='scrum'
	//	   search all employees working on a specific task		//       /tasks?name=

	//http.HandleFunc("/employees", EMP)
	//	http.HandleFunc("/tasks", TAS)
	//	http.HandleFunc("/whois", WHO)
	//http.HandleFunc("/newEmployee", empPOST)
	//	http.HandleFunc("/newTask", POSTHandler)

	//log.Fatal(http.ListenAndServe(":8080", nil))
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/employees", returnAllEmployees)
	myRouter.HandleFunc("/employees/uuid/{uniqid}", returnSingleByUUID)
	myRouter.HandleFunc("/employees/empid/{empid}", returnSingleArticle)
	myRouter.HandleFunc("/tasks/name/{title}", returnSingleTask)

	//myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	//myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	//myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

//func main() {
//	Articles = []Article{
//		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
//		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
//	}

//}

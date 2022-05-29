package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func OpenConnection() *sql.DB {

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func EMP(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM employees")
	if err != nil {
		log.Fatal(err)
	}

	var employees []Employee

	for rows.Next() {
		var employee Employee
		rows.Scan(&employee.uniqid, &employee.empid, &employee.fname, &employee.lname)
		employees = append(employees, employee)
	}

	empBytes, _ := json.MarshalIndent(employees, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(empBytes)

	defer rows.Close()
	defer db.Close()
}

func TAS(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
	}

	var todos []Task

	for rows.Next() {
		var task Task
		rows.Scan(&task.unitid, &task.taskid, &task.assignedto, &task.title, &task.privacy)
		todos = append(todos, task)
	}

	todosBytes, _ := json.MarshalIndent(todos, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(todosBytes)

	defer rows.Close()
	defer db.Close()
}

func WHO(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	var e Whois
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `SELECT * FROM employees INNER JOIN tasks ON empid=ANY(assignedto) WHERE title IN ($1)`
	rows, err := db.Query(sqlStatement, e.title)
	if err != nil {
		log.Fatal(err)
	}

	var assigned []Whois

	for rows.Next() {
		var whois Whois
		rows.Scan(&whois.uniqid, &whois.empid, &whois.fname, &whois.lname, &whois.unitid, &whois.taskid, &whois.assignedto, &whois.title, &whois.privacy)
		assigned = append(assigned, whois)
	}

	whoBytes, _ := json.MarshalIndent(assigned, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(whoBytes)

	defer rows.Close()
	defer db.Close()
}

func empPOST(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var e Employee
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO employees (empid, fname, lname) VALUES ($1, $2,$3)`
	_, err = db.Exec(sqlStatement, e.empid, e.fname, e.lname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func tasPOST(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var t Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO tasks (assignedto, title, privacy) VALUES ($1, $2,$3)`
	_, err = db.Exec(sqlStatement, t.assignedto, t.title, t.privacy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

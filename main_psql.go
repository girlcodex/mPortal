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
		rows.Scan(&employee.Uniqid, &employee.Empid, &employee.Fname, &employee.Lname)
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
		rows.Scan(&task.Unitid, &task.Taskid, &task.Assignedto, &task.Title, &task.Privacy)
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
	rows, err := db.Query(sqlStatement, e.Title)
	if err != nil {
		log.Fatal(err)
	}

	var assigned []Whois

	for rows.Next() {
		var whois Whois
		rows.Scan(&whois.Uniqid, &whois.Empid, &whois.Fname, &whois.Lname, &whois.Unitid, &whois.Taskid, &whois.Assignedto, &whois.Title, &whois.Privacy)
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
	_, err = db.Exec(sqlStatement, e.Empid, e.Fname, e.Lname)
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
	_, err = db.Exec(sqlStatement, t.Assignedto, t.Title, t.Privacy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

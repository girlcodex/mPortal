package main

import (
	"log"
)

func mainPSQL() {

	// check for data presence
	db := OpenConnection()

	_, err := db.Query("SELECT * FROM employees")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	//if no data present in psql, install sample data

	// vv--- this needs to be inside an iterator that takes the json and parses it
	//	sqlStatement := "INSERT INTO tasks (assignedto, title, privacy) VALUES ($1, $2,$3)"
	//	_, err = db.Exec(sqlStatement, t.assignedto, t.title, t.privacy)
	//	if err != nil {
	//		panic(err)
	//	}
	///	defer db.Close()
	//if data is present, sync to ES
}

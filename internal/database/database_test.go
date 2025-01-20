package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"

	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	_ "github.com/mattn/go-sqlite3"
)

// public function for initializatin the database
func TestDatabaseInit(t *testing.T) {
	db, err := sql.Open("sqlite3", "./testdatabase.db")
	if err != nil {
		log.Fatalln(err)
	}
	// defer db.Close()

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS testdb (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = db.Prepare("INSERT INTO testdb (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Enrique", "Romero")

	rows, _ := db.Query("SELECT id, firstname, lastname FROM testdb")

	var id int
	var firstname string
	var lastname string

	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + " : " + firstname + " " + lastname)
	}

}

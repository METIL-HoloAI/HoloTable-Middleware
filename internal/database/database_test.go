package database_test

import (
	"database/sql"
	"fmt"
	"os"

	"strconv"
	"testing"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

// public function for initializing the database
func TestDatabaseInit(t *testing.T) {
	config.LoadYaml()
	if err := os.MkdirAll(config.General.DataDir, os.ModePerm); err != nil {
		t.Fatal("Failed to create data directory:", err)
	}
	db, err := sql.Open("sqlite3", config.General.DataDir+"test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS testdb (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		t.Fatal(err)
	}
	statement, err = db.Prepare("INSERT INTO testdb (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = statement.Exec("Enrique", "Romero")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("SELECT id, firstname, lastname FROM testdb")
	if err != nil {
		t.Fatal(err)
	}

	var id int
	var firstname string
	var lastname string

	for rows.Next() {
		err := rows.Scan(&id, &firstname, &lastname)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(strconv.Itoa(id) + " : " + firstname + " " + lastname)
	}

}

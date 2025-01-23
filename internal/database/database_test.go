package database

import (
	"database/sql"
	"fmt"
	"os"

	// "log"
	"strconv"
	"testing"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	_ "github.com/mattn/go-sqlite3"
)

// public function for initializing the database
func TestDatabaseInit(t *testing.T) {
	// Load yaml
	settings, err := configloader.GetGeneral()
	if err != nil {
		t.Fatal("Error loading general settings")
		t.Fatal(err)
		return
	}

	if err := os.MkdirAll(settings.DataDir, os.ModePerm); err != nil {
		t.Fatal("Failed to create data directory:", err)
	}
	db, err := sql.Open("sqlite3", settings.DataDir+"test.db")
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

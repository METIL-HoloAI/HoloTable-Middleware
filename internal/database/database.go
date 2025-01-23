package database

import (
	"database/sql"
	"fmt"
	"os"

	"log"
	"strconv"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// public function for initializing the database
func Init() {
	// Load yaml
	settings, err := configloader.GetGeneral()
	if err != nil {
		log.Fatal("Error loading general settings")
		log.Fatal(err)
		return
	}

	db, err = sql.Open("sqlite3", settings.DataDir+"filelocations.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO:
	// need table for each type of thing (image, video, gifs, 3dmodels)
	// store file location
	// store filename

	fileTypes := []string{"image", "video", "gif", "model"}
	// iterates through the four file types and creates a table for each
	// each table contains id and filepath
	for i := 0; i < len(fileTypes); i++ {
		statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS" + fileTypes[i] + "(id INTEGER PRIMARY KEY, filepath TEXT)")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
		statement, err = db.Prepare("INSERT INTO" + fileTypes[i] + "(filepath) VALUES (?)")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec("some filepath")

		rows, err := db.Query("SELECT id, filepath FROM" + fileTypes[i])
		if err != nil {
			log.Fatal(err)
		}

		var id int
		var filepath string

		for rows.Next() {
			rows.Scan(&id, &filepath)
			fmt.Println(strconv.Itoa(id) + " : " + filepath + " ")
		}
	}
}

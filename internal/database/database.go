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

	if err := os.MkdirAll(settings.DataDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}
	db, err = sql.Open("sqlite3", settings.DataDir+"filelocations.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
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
		statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS " + fileTypes[i] + " (id INTEGER PRIMARY KEY, filepath TEXT)")
		if err != nil {
			log.Fatalf("Failed to prepare CREATE TABLE statement for '%s': %v", fileTypes[i], err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatalf("Failed to execute CREATE TABLE statement for '%s': %v", fileTypes[i], err)
		}
		statement, err = db.Prepare("INSERT INTO " + fileTypes[i] + " (filepath) VALUES (?)")
		if err != nil {
			log.Fatalf("Failed to prepare INSERT statement for '%s': %v", fileTypes[i], err)
		}
		_, err = statement.Exec("some filepath")
		if err != nil {
			log.Fatalf("Failed to execute INSERT statement for '%s': %v", fileTypes[i], err)
		}

		rows, err := db.Query("SELECT id, filepath FROM " + fileTypes[i])
		if err != nil {
			log.Fatalf("Failed to query SELECT statement for '%s': %v", fileTypes[i], err)
		}

		var id int
		var filepath string

		for rows.Next() {
			if err := rows.Scan(&id, &filepath); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			if _, err := fmt.Println(strconv.Itoa(id) + " : " + filepath + " "); err != nil {
				log.Printf("Error writing output: %v", err)
			}
		}
	}
}

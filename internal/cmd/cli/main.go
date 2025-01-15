package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load yaml
	settings, err := configloader.GetGeneral()
	if err != nil {
		fmt.Println("Error loading general settings")
		fmt.Println(err)
		return
	}

	// Load database
	db, err := sql.Open("sqlite3", settings.DataDir+"/db.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Check how user wants to listen for input
	// and start that listener
	if settings.Listener == "mic" {
		fmt.Println("Microphone Listener")
	} else if settings.Listener == "text" {
		fmt.Println("Text Listener")
	} else {
		fmt.Println("Invalid listener option in general.yaml")
	}

}

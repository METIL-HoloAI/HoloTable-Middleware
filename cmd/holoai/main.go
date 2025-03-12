package main

import (
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// config.LoadYaml()

	// if err := os.MkdirAll(config.General.DataDir, os.ModePerm); err != nil {
	// 	log.Fatal("Failed to create data directory:", err)
	// }

	// db, err := sql.Open("sqlite3", config.General.DataDir+"database.db?_mode=shared&_journal_mode=WAL")
	// if err != nil {
	// 	log.Fatal("Failed to open database:", err)
	// }
	// defer db.Close()

	// database.Init(db)

	// Start WebSocket server // IMPORTANT TO ADD THIS
	go unityserver.StartWebSocketServer()
	<-unityserver.ClientReady

	// unityserver.ExportAsset(fileName, extension, filePath) // HOW TO call my function

	// Check how user wants to listen for input
	// and start that listener
	// if config.General.Listener == "mic" {
	// 	fmt.Println("Microphone Listener")
	// } else if config.General.Listener == "text" {
	// 	listeners.StartTextListener()
	// } else {
	// 	fmt.Println("Invalid listener option in general.yaml")
	// }
}

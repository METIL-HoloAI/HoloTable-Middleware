package database_test

import (
	"log"

	"testing"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

// public function for initializing the database
func TestDatabaseInit(t *testing.T) {
	config.LoadYaml()
	database.Init()

	err := database.Insert("image", "test", "test")
	if err != nil {
		t.Fatal("Failed to insert into image, ", err)
	}

	path, err := database.GetPathByFilename("image", "test")
	if err != nil {
		t.Fatal("Failed to get test back, ", err)
	}
	t.Log(path)

	items, err := database.ListAllFilenames("image")
	if err != nil {
		t.Fatal("Failed to list all filenames in images, ", err)
	}
	for _, item := range items {
		t.Log(item)
	}

	err = database.DeleteRecordByFilename("image", "test")
	if err != nil {
		log.Fatal("Failed to remove test record, ", err)
	}
}

package database_test

import (
	"database/sql"
	"os"

	"testing"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// public function for initializing the database
func TestDatabaseInit(t *testing.T) {
	if err := os.MkdirAll("./testdb/", os.ModePerm); err != nil {
		t.Fatal("Failed to create data directory:", err)
	}

	db, err := sql.Open("sqlite3", config.General.DataDir+"test.db")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	database.Init(db)

	err = database.Insert("image", "test", "test")
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
		logrus.Fatal("Failed to remove test record, ", err)
	}
}

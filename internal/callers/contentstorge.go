package callers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config/structs"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
)

var General structs.GeneralSettings

const filePerm = 0644

// ContentStorage saves the content to the specified path in local storage
func ContentStorage(fileType, filename string, content []byte) error {
	// Map file types to database table names
	tableMap := map[string]string{
		"image": "image",
		"3d":    "model",
		"gif":   "gif",
		"video": "video",
	}

	// Get the corresponding database table name
	tableName, ok := tableMap[fileType]
	if !ok {
		return fmt.Errorf("invalid file type: %s", fileType)
	}

	// Create the subdirectory path
	directory := filepath.Join(General.DataDir, tableName)

	// Ensure the directory exists
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the full file path
	filePath := filepath.Join(directory, filename)

	// Write the content to the file first before inserting into the database
	if err := os.WriteFile(filePath, content, filePerm); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	// Insert into the database after a successful write
	if err := database.Insert(tableName, filename, filePath); err != nil {
		return fmt.Errorf("failed to insert record into database: %v", err)
	}

	return nil
}

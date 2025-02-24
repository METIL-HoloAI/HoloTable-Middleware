package callers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config/structs"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
)

var General structs.GeneralSettings

const filePerm = 0644

// ContentStorage saves the content to local storage under a subdirectory based on the file type.
// If the provided content represents a URL (i.e. starts with "http://" or "https://"),
// the function downloads the file from that URL before storing it.
func ContentStorage(fileType, filename string, content []byte) error {
	// Map file types to database table names.
	tableMap := map[string]string{
		"image": "image",
		"3d":    "model",
		"gif":   "gif",
		"video": "video",
	}

	// Get the corresponding database table name.
	tableName, ok := tableMap[fileType]
	if !ok {
		return fmt.Errorf("invalid file type: %s", fileType)
	}

	// Check if the content is a URL by converting it to a string.
	contentStr := strings.TrimSpace(string(content))
	if strings.HasPrefix(contentStr, "http://") || strings.HasPrefix(contentStr, "https://") {
		// Download the file from the URL.
		resp, err := http.Get(contentStr)
		if err != nil {
			return fmt.Errorf("failed to download content from URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download content: received status code %d", resp.StatusCode)
		}

		downloadedContent, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read downloaded content: %v", err)
		}
		// Replace content with the downloaded data.
		content = downloadedContent
	}

	// Create the subdirectory path.
	directory := filepath.Join(General.DataDir, tableName)
	// Ensure the directory exists.
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the full file path.
	filePath := filepath.Join(directory, filename)
	// Write the content to the file.
	if err := os.WriteFile(filePath, content, filePerm); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	// Insert a record into the database.
	if err := database.Insert(tableName, filename, filePath); err != nil {
		return fmt.Errorf("failed to insert record into database: %v", err)
	}

	return nil
}

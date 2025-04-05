package callers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config/structs"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/google/uuid"
)

var General structs.GeneralSettings

const filePerm = 0644

// ContentStorage saves the content to local storage under a subdirectory based on the file type.
// If the provided content represents a URL (i.e. format == "url"), the function downloads the file from that URL before storing it.
// It assumes that the provided filename already includes the proper file extension.
func ContentStorage(fileType, format, fileExtention string, content []byte) ([]byte, string, string, error) {
	// Map file types to database table names.
	tableMap := map[string]string{
		"image": "image",
		"model": "model",
		"gif":   "gif",
		"video": "video",
	}

	// Get the corresponding database table name.
	tableName, ok := tableMap[fileType]
	if !ok {
		return nil, "", "", fmt.Errorf("invalid file type: %s", fileType)
	}

	// Combine fileID and fileExtention into a single file name.
	fileID := uuid.New().String()
	fileName := ""
	if fileExtention != "" {
		fileName = fmt.Sprintf("%s.%s", fileID, fileExtention)
	}

	// If the format indicates that the content is a URL, download the file.
	if format == "url" {
		var err error
		content, err = downloadContent(string(content))
		if err != nil {
			return nil, "", "", err
		}
	}

	// Create the subdirectory path.
	directory := filepath.Join(config.General.DataDir, "/content", tableName)
	// Ensure the directory exists.
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, "", "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the full file path.
	filePath := filepath.Join(directory, fileName)
	// Write the content to the file.
	if err := os.WriteFile(filePath, content, filePerm); err != nil {
		return nil, "", "", fmt.Errorf("failed to write file: %v", err)
	}

	// Insert a record into the database.
	if err := database.Insert(tableName, fileName, filePath); err != nil {
		return nil, "", "", fmt.Errorf("failed to insert record into database: %v", err)
	}

	// filePath = "my_file.txt" // Could be a relative or absolute path

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return nil, "", "", err
	}

	fmt.Println("Absolute path:", absPath)

	return content, absPath, fileID, nil
}

// downloadContent downloads the content from the given URL and returns the downloaded data.
func downloadContent(urlStr string) ([]byte, error) {
	resp, err := http.Get(strings.TrimSpace(urlStr))
	if err != nil {
		return nil, fmt.Errorf("failed to download content from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download content: received status code %d", resp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read downloaded content: %v", err)
	}

	return downloadedContent, nil
}

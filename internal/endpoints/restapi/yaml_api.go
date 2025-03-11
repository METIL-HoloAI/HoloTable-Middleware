package restapi

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/gorilla/mux"
)

// https://pkg.go.dev/github.com/gorilla/mux
// https://medium.com/better-programming/building-a-simple-rest-api-in-go-with-gorilla-mux-892ceb128c6f
// send text

func YamlApi() {
	router := mux.NewRouter()
	router.HandleFunc("/config/{name}", getYamlHandler).Methods("GET")
	router.HandleFunc("/config/{name}", putYamlHandler).Methods("PUT")
	router.HandleFunc("/database/list", useListAllFilenames).Methods("PUT")

	// start serv
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getYamlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	yamlName := vars["name"]
	yamlPath := "../../config/" + yamlName + ".yaml"

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		yamlPath = "../../config/contentgen/" + yamlName + ".yaml"
		data2, err := os.ReadFile(yamlPath)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		data = data2
	}

	dataString := string(data)

	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(dataString))
	if err != nil {
		log.Fatal("Failed to write response:", err)
	}
}

func putYamlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	yamlName := vars["name"]
	yamlPaths := []string{
		"../../config/" + yamlName + ".yaml",
		"../../config/contentgen/" + yamlName + ".yaml",
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingPath string
	for _, path := range yamlPaths {
		if _, err := os.Stat(path); err == nil {
			existingPath = path
			break
		}
	}

	if existingPath == "" {
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}

	err = os.WriteFile(existingPath, body, 0644)
	if err != nil {
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("YAML updated"))
	if err != nil {
		log.Fatal("Failed to write response:", err)
	}
}

// create routes thru database
// Allow users to call and use ListAllFilenames thru an API
func useListAllFilenames(w http.ResponseWriter, r *http.Request) {
	tables := [4]string{"image", "video", "gif", "model"}
	var allFilenames []string

	for i := 0; i < 4; i++ {
		cache, err := database.ListAllFilenames(tables[i])
		if err != nil {
			http.Error(w, "Failed to load database table"+tables[i], http.StatusBadRequest)
			return
		}
		allFilenames = append(allFilenames, cache...)
	}
	// converts string array to byte slice
	responseBytes := []byte(strings.Join(allFilenames, "\n"))

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write(responseBytes)
	if err != nil {
		log.Fatal("Failed to write response:", err)
	}
}

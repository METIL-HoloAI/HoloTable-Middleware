package restapi

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/gorilla/mux"
)

// https://pkg.go.dev/github.com/gorilla/mux
// https://medium.com/better-programming/building-a-simple-rest-api-in-go-with-gorilla-mux-892ceb128c6f
// send text

func enableCORS(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			// If it's a preflight OPTIONS request, exit here.
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func StartRestAPI() {
	router := mux.NewRouter()
	// Enable CORS for the router
	enableCORS(router)

	// New endpoint for keyword
	router.HandleFunc("/config/keyword", getKeywordHandler).Methods("GET", "OPTIONS")

	// ({name: .*} essentially is a "catch-all route" meaning it will catch the rest of the route after "/config"
	// This ensures that config files in lower directories can still be fetched, i.e. "/config/contentgen_yamls" is all captured
	router.HandleFunc("/config/{name:.*}", getYamlHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/config/{name:.*}", putYamlHandler).Methods("PUT", "OPTIONS")
	router.HandleFunc("/database/list", useListAllFilenames).Methods("PUT", "OPTIONS")

	// start serv
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getKeywordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// Return the keyword from the SpeechToText configuration
	_, err := w.Write([]byte(config.SpeechToText.Keyword))
	if err != nil {
		log.Fatal("Failed to write response:", err)
	}
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

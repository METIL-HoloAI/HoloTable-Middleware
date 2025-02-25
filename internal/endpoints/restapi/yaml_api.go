package restapi

import (
	"log"
	"net/http"
	"os"
	"io"
	"strings"
	"github.com/gorilla/mux"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
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
	w.Write([]byte(dataString))
}

func putYamlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	yamlName := vars["name"]
	yamlPath := "../../config/" + yamlName + ".yaml"

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err = os.WriteFile(yamlPath, body, 0)
    if err != nil {
		yamlPath = "../../config/contentgen/" + yamlName + ".yaml"
		err2 := os.WriteFile(yamlPath, body, 0)
		if err2 != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("YAML updated"))
}

// create routes thru database
// Allow users to call and use ListAllFilenames thru an API 
func useListAllFilenames(w http.ResponseWriter, r *http.Request) {
	tables := [4]string{"image", "video", "gif", "model"}
	var allFilenames []string

	for i := 0; i < 4; i++ {
		cache, err := database.ListAllFilenames(tables[i])
		if err != nil {
			http.Error(w, "Failed to load database table" + tables[i], http.StatusBadRequest)
			return
		}
		allFilenames = append(allFilenames, cache...)
	}
	// converts string array to byte slice
	responseBytes := []byte(strings.Join(allFilenames, "\n"))

	w.Header().Set("Content-Type", "text/plain")
	w.Write(responseBytes)
}
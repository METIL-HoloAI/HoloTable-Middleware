package restapi

import (
	"log"
	"net/http"
	"os"
	"io"
	"github.com/gorilla/mux"
)

// https://pkg.go.dev/github.com/gorilla/mux
// https://medium.com/better-programming/building-a-simple-rest-api-in-go-with-gorilla-mux-892ceb128c6f
// send text
// create routes thru database

func YamlApi() {
	router := mux.NewRouter()
	router.HandleFunc("/config/{name}", getYamlHandler).Methods("GET")
	router.HandleFunc("/config/{name}", putYamlHandler).Methods("PUT")

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
        http.Error(w, "Failed to write file", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("YAML updated"))
}

package restapi

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// https://pkg.go.dev/github.com/gorilla/mux
// https://medium.com/better-programming/building-a-simple-rest-api-in-go-with-gorilla-mux-892ceb128c6f
// send text
// create routes thru database

func YamlApi() {
	router := mux.NewRouter()
	router.HandleFunc("/config/{name}", getYamlHandler).Methods("GET")
	//router.HandleFunc("/config/{name}", putYamlHandler).Methods("PUT")

	// start serv
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getYamlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	yamlName := vars["name"]
	yamlPath := yamlName + ".yaml"

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Write(data)
}

// func putYamlHandler(w http.ResponseWriter, r *http.Request) {

// }

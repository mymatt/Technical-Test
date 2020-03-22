package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Myapp json return for /version endpoint
type Myapp struct {
	App []Appdetails `json:"myapplication"`
}

type Appdetails struct {
	Version       string `json:"version"`
	Description   string `json:"description"`
	Lastcommitsha string `json:"lastcommitsha"`
}

// endpoint / return
func getHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}

// endpoint /version return
func getVersion(w http.ResponseWriter, r *http.Request) {

	version, exists := os.LookupEnv("VERS")
	if !exists {
		panic(exists)
	}

	description, exists := os.LookupEnv("DESC")
	if !exists {
		panic(exists)
	}

	shacommit, exists := os.LookupEnv("SHA")
	if !exists {
		panic(exists)
	}

	myapp := &Myapp{
		App: []Appdetails{
			{Version: version, Description: description, Lastcommitsha: shacommit},
		},
	}

	appJSON, err := json.Marshal(myapp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(appJSON)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", getHome).Methods(http.MethodGet)
	router.HandleFunc("/version", getVersion).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", router))
}

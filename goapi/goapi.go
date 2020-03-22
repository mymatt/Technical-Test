package main

import "net/http"

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

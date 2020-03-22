package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Results struct {
	Myapplication []Version
}

type Version struct {
	Version       string
	Description   string
	Lastcommitsha string
}

var test_hello = "Hello World!"
var test_version = "1.0.0"
var test_description = "anz technical challenge"
var test_lastcommitsha = "ttn5int5ln34r34"

func TestHome(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getHome)

	handler.ServeHTTP(rec, req)

	// Test http status code
	if status := rec.Code; status == http.StatusOK {
		fmt.Printf("-> Test getHome() Status Code Passed \nReceived correct status code: %v \n", status)
	} else {
		t.Errorf("-> Received status code: %v expected status code was: %v",
			status, http.StatusOK)
	}

	if status := rec.Body.String(); status == test_hello {
		fmt.Printf("-> Test getHome() Body Passed \nReceived correct body: %v \n", status)
	} else {
		t.Errorf("Received body: %v expected body was %v",
			status, test_hello)
	}
}

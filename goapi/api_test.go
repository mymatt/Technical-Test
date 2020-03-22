package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func TestVersion(t *testing.T) {
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getVersion)

	handler.ServeHTTP(rec, req)

	// Test http status code
	if status := rec.Code; status == http.StatusOK {
		fmt.Printf("-> Test /version endpoint Passed \nReceived correct status code: %v \n", status)
	} else {
		t.Errorf("-> Received status code: %v expected status code was %v",
			status, http.StatusOK)
	}

	res := &Results{}
	er := json.Unmarshal([]byte(rec.Body.String()), res)
	if er != nil {
		log.Fatal(er)
	}

	// Test /version endpoint: Version
	if version := res.Myapplication[0].Version; version == test_version {
		fmt.Printf("-> Test Version Passed \nReceived correct version: %v \n", version)
	} else {
		t.Errorf("-> Received version: %v expected version was %v",
			version, test_version)
	}

	// Test /version endpoint: Description
	if description := res.Myapplication[0].Description; description == test_description {
		fmt.Printf("-> Test Description Passed \nReceived correct description: %v \n", description)
	} else {
		t.Errorf("-> Received description: %v expected description was %v",
			description, test_description)
	}

	// Test /version endpoint: Lastcommitsha
	// if lastcommitsha := res.Myapplication[0].Lastcommitsha; lastcommitsha == test_lastcommitsha {
	// 	fmt.Printf("-> Test Lastcommitsha Passed \nReceived correct lastcommitsha: %v \n", lastcommitsha)
	// } else {
	// 	t.Errorf("-> Received Lastcommitsha: %v expected Lastcommitsha was %v",
	// 		lastcommitsha, test_lastcommitsha)
	// }

}

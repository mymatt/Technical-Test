package main

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

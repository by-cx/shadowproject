package main

import (
	"log"
	"net/http"
	"shadowproject/common"
	"shadowproject/docker"
)

// Test task
var MyTask = common.Task{
	Driver:  &docker.DockerDriver{},
	UUID:    "bainimiepaevaichaoloeloneisieshu",
	Domains: []string{"localhost"},
	Image:   "creckx/testimage",
	Command: []string{"/srv/testtask"},
}

func init() {
	// Prepare the environment
	driver := docker.DockerDriver{}
	driver.Clear()
}

func main() {
	// Set up the reverse proxy
	http.HandleFunc("/", ReverseProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

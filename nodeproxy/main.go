package main

import (
	"log"
	"net/http"
	"shadowproject/docker"
)

var config NodeProxyConfig

func init() {
	// Handle config
	ProcessEnvironmentVariables(&config)
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

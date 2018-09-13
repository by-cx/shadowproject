package main

import (
	"log"
	"net/http"
	"shadowproject/docker"
	"shadowproject/master/client"
)

var config NodeProxyConfig
var shadowClient client.ShadowMasterClientInterface

func init() {

}

func main() {
	// Handle config
	ProcessEnvironmentVariables(&config)

	// Shadow client to connect to master server
	shadowClient = &client.ShadowMasterClient{
		Host: config.MasterHost,
		Port: config.MasterPort,
	}

	// Prepare the environment
	driver := docker.DockerDriver{}
	driver.Clear()

	// Set up the reverse proxy
	http.HandleFunc("/", ReverseProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

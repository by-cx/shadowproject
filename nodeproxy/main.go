package main

import (
	"log"
	"net/http"
	"shadowproject/docker"
	"shadowproject/master/client"
	"strconv"
)

var config NodeProxyConfig
var shadowClient client.ShadowMasterClientInterface
var dockerDriver docker.ContainerDriverInterface

// TODO: kill containers after configured amount of time
// TODO: remove the TaskCache for now, check for existence of the container

func main() {
	// Handle config
	ProcessEnvironmentVariables(&config)

	// Shadow client to connect to master server
	shadowClient = &client.ShadowMasterClient{
		Host:  config.MasterHost,
		Port:  config.MasterPort,
		Proto: config.MasterProto,
	}

	// Prepare the environment
	dockerDriver = &docker.DockerDriver{}
	dockerDriver.Clear()

	// Set up the reverse proxy
	http.HandleFunc("/", ReverseProxyHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

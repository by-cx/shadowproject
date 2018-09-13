package main

import (
	"log"
	"net/http"
	"shadowproject/docker"
	"shadowproject/master/client"
	"strconv"
	"time"
)

var config NodeProxyConfig
var shadowClient client.ShadowMasterClientInterface
var dockerDriver docker.ContainerDriverInterface
var LastRequestMap = make(map[string]int64) // Map where key is time of the last request and value is TaskUUID

// TODO: kill containers after configured amount of time

// After this amount of seconds without any request, the container is killed
const KILL_AFTER = 10

func ContainerCleaner() {
	now := time.Now().Unix()

	for taskUUID, lastRequestTime := range LastRequestMap {
		log.Println("TaskUUID clenaer:", taskUUID)
		if now-lastRequestTime > KILL_AFTER {
			task, err := shadowClient.GetTask(taskUUID)
			if err != nil {
				log.Println("Container cleaner error:", err)
				continue
			}
			log.Println("Removing containers for ", taskUUID)
			log.Println("Task:", task)
			task.Driver = dockerDriver
			task.DestroyAll()
		}
	}
}

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
	log.Println("Ready to accept connections ..")

	// Container cleaner
	go func() {
		log.Println("Starting container cleaner ..")
		for {
			ContainerCleaner()
			time.Sleep(time.Second * (KILL_AFTER + 5))
		}
		log.Println("Stopping container cleaner ..")
	}()

	// Set up the reverse proxy
	http.HandleFunc("/", ReverseProxyHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

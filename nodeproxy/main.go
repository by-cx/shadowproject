package main

import (
	"log"
	"net/http"
	"shadowproject/common"
	"shadowproject/docker"
)

var MyTask = common.Task{
	Driver:  &docker.DockerDriver{},
	UUID:    "bainimiepaevaichaoloeloneisieshu",
	Domains: []string{"localhost"},
	Image:   "creckx/testimage",
	Command: []string{"/srv/testtask"},
}

func main() {
	http.HandleFunc("/", ReverseProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

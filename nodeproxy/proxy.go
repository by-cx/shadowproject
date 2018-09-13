package main

import (
	"io"
	"log"
	"net/http"
	"shadowproject/common"
	shadowerrors "shadowproject/common/errors"
	"strconv"
	"time"
)

var TaskCache = make(map[string]*common.Task)

// Get task from cache if it's possible otherwise use
// TODO: check regularly changes in cached tasks, remove all containers in case of change and create new ones
func GetTaskByDomain(domain string) *common.Task {
	var err error
	var task *common.Task

	if val, ok := TaskCache[domain]; ok {
		return val
	} else {
		task, err = shadowClient.GetTaskByDomain(domain)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "backend connection error",
			})
		}
		TaskCache[domain] = task
	}

	return task
}

// Return HOST:PORT pair to connect to for the specified domain.
func FindContainer(domain string) (string, error) {
	var containerUUID string
	var err error

	task := GetTaskByDomain(domain)

	if len(task.Containers) > 0 {
		containerUUID = task.Containers[0]
	} else {
		containerUUID, err = task.AddContainer()
		if err != nil {
			return "", err
		}
	}

	port, err := task.Driver.GetPort(containerUUID)
	if err != nil {
		return "", err
	}

	return "localhost:" + strconv.Itoa(port), nil
}

// This takes incoming request, updated URL struct and sends it to the backend.
// Then the response is streamed back including headers.
func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	var remoteRequest http.Request
	remoteRequest = *r

	// Find the destination
	domain, err := FindContainer(r.URL.Hostname())
	if err != nil {
		panic(err)
	}

	log.Println(r.Host, "directed to "+domain)

	// Copy request and update some values
	remoteRequest.RequestURI = ""
	remoteRequest.Host = domain
	remoteRequest.URL.Host = domain
	remoteRequest.URL.Scheme = "http"

	// Do the remote request
	httpClient := http.Client{
		Timeout: time.Second * PROXY_TIMEOUT,
	}
	response, err := httpClient.Do(&remoteRequest)
	if err != nil {
		panic(err)
	}

	// Pass headers
	for key, values := range response.Header {
		for i, value := range values {
			if i == 0 {
				w.Header().Set(key, value)
			} else {
				w.Header().Add(key, value)
			}
		}
	}

	// Set status code
	w.WriteHeader(response.StatusCode)

	// Stream the body
	p := make([]byte, PROXY_BUFFER_SIZE)
	for {
		n, err := response.Body.Read(p)

		// If we have data send it
		if n > 0 {
			w.Write(p[0:n])
		}

		// If it's end of the stream bye bye
		if err == io.EOF {
			break
		}
	}
}

package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"shadowproject/common"
	shadowerrors "shadowproject/common/errors"
	"strconv"
	"time"
)

// Gets info about the domain from the master server
func GetTaskByDomain(domain string) *common.Task {
	var err error
	var task *common.Task

	task, err = shadowClient.GetTaskByDomain(domain)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "backend connection error",
		})
	}

	return task
}

// Return HOST:PORT pair to connect to for the specified domain.
func FindContainer(domain string) (string, error) {
	var containerUUID string

	// Get the task
	task := GetTaskByDomain(domain)
	task.ContainerDriver = dockerDriver
	if task.VolumeType == common.VolumeTypeS3 {
		task.VolumeDriver = S3VolumeDriver
	} else {
		panic(shadowerrors.ShadowError{
			VisibleMessage: "unknown volume driver",
			Origin:         errors.New("unknown volume driver"),
		})
	}

	// Mark time of this request
	LastRequestMap[task.UUID] = time.Now().Unix()

	// Find containers and decide whether new container is needed
	containerUUIDs := task.ContainerDriver.IsExist(task.UUID)

	if len(containerUUIDs) > 0 {
		containerUUID = containerUUIDs[rand.Int()%len(containerUUIDs)]
	} else {
		containerUUID = task.AddContainer()
	}

	// Get the port where to point the request
	port := task.ContainerDriver.GetPort(containerUUID)

	// Return the destination
	return config.ProxyTarget + ":" + strconv.Itoa(port), nil
}

// Returns json string with server status
func getJSONedHealth() []byte {
	status := HealthCheck()
	content, err := json.Marshal(status)

	if err != nil {
		panic(shadowerrors.ShadowError{
			VisibleMessage: "health check error",
			Origin:         errors.New("unknown volume driver"),
		})
	}

	return content
}

// This takes incoming request, updated URL struct and sends it to the backend.
// Then the response is streamed back including headers.
func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	// In case of failure, print the error message so user can see what is happening
	defer func() {
		if r := recover(); r != nil {
			if config.Debug {
				debug.PrintStack()
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(r.(error).Error()))
		}
	}()

	// Health check
	if r.URL.Path == "/server-health-check" {
		w.WriteHeader(http.StatusOK)
		w.Write(getJSONedHealth())
		return
	}

	// Processing the reverse proxy request
	var remoteRequest http.Request
	remoteRequest = *r

	// Find the destination
	domain, err := FindContainer(r.Host)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "backend connection error",
		})
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

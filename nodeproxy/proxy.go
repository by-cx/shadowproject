package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

// This takes incoming request, updated URL struct and sends it to the backend.
// Then the response is streamed back including headers.
func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	var remoteRequest http.Request
	remoteRequest = *r

	var domain = "localhost:32770"

	log.Println(r.Host, "directed to "+domain)

	// Copy request and update some values
	remoteRequest.RequestURI = ""
	remoteRequest.Host = domain
	remoteRequest.URL.Host = domain
	remoteRequest.URL.Scheme = "http"

	// Do the remote request
	client := http.Client{
		Timeout: time.Second * PROXY_TIMEOUT,
	}
	response, err := client.Do(&remoteRequest)
	if err != nil {
		log.Fatal(err)
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

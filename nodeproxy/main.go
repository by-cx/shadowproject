package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", ReverseProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

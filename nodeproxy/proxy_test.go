package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTaskByDomain(t *testing.T) {
	shadowClient = &MockShadowMasterClient{}
	dockerDriver = &MockDockerDriver{}
	ProcessEnvironmentVariables(&config)

	task := GetTaskByDomain("localhost")
	assert.Equal(t, task.Domains, []string{"localhost"})
}

func TestFindContainer(t *testing.T) {
	shadowClient = &MockShadowMasterClient{}
	dockerDriver = &MockDockerDriver{}
	S3VolumeDriver = &MockS3VolumeDriver{}
	ProcessEnvironmentVariables(&config)

	containerId, err := FindContainer("localhost")

	assert.Equal(t, "localhost:80", containerId)
	assert.Nil(t, err)
}

func TestReverseProxyHandler(t *testing.T) {
	shadowClient = &MockShadowMasterClient{}
	dockerDriver = &MockDockerDriver{}
	ProcessEnvironmentVariables(&config)
	config.ProxyTarget = "ifconfig.co"

	req, err := http.NewRequest("GET", "/", nil)
	req.Host = "ifconfig.co"
	req.Header.Add("Host", "ifconfig.co")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReverseProxyHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTaskByDomain(t *testing.T) {
	shadowClient = &MockShadowMasterClient{}
	dockerDriver = &MockDockerDriver{}

	task := GetTaskByDomain("localhost")
	assert.Equal(t, task.Domains, []string{"localhost"})
}

func TestFindContainer(t *testing.T) {
	shadowClient = &MockShadowMasterClient{}
	dockerDriver = &MockDockerDriver{}

	containerId, err := FindContainer("localhost")

	assert.Equal(t, "localhost:32000", containerId)
	assert.Nil(t, err)
}

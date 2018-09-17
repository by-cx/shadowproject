package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	// Set up
	dockerDriver = &MockDockerDriver{}

	// the test
	var status Status

	status = *HealthCheck()

	assert.True(t, status.CPUUtilization > 0)
	assert.True(t, status.Memory > 0)
	assert.True(t, status.ContainerBackend)
	assert.False(t, status.Critical)
}

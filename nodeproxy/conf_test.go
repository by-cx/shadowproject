package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProcessEnvironmentVariables(t *testing.T) {
	os.Setenv("HOST", "foo.local")
	os.Setenv("PORT", "1234")

	var config NodeProxyConfig
	ProcessEnvironmentVariables(&config)

	assert.Equal(t, config.MasterPort, 1234)
	assert.Equal(t, config.MasterHost, "foo.local")

	os.Remove("HOST")
	os.Remove("PORT")
}

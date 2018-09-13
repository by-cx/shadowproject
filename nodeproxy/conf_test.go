package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProcessEnvironmentVariables(t *testing.T) {
	os.Setenv("MASTER_HOST", "foo.local")
	os.Setenv("MASTER_PORT", "1234")
	os.Setenv("MASTER_PROTO", "xyz")
	os.Setenv("PORT", "3214")

	var config NodeProxyConfig
	ProcessEnvironmentVariables(&config)

	assert.Equal(t, config.MasterPort, 1234)
	assert.Equal(t, config.MasterHost, "foo.local")
	assert.Equal(t, config.MasterProto, "xyz")
	assert.Equal(t, config.Port, 3214)

	os.Remove("MASTER_HOST")
	os.Remove("MASTER_PORT")
	os.Remove("MASTER_PROTO")
	os.Remove("PORT")
}

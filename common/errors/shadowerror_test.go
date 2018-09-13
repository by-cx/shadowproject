package errors

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShadowError_Error(t *testing.T) {
	err := &ShadowError{
		VisibleMessage: "test visible message",
		Origin:         errors.New("origin error"),
	}

	assert.Equal(t, err.Error(), "origin error")
	assert.Equal(t, err.Message(), "test visible message")
}

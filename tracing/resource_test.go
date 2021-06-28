package tracing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")
	initConfig()
	res, _ := newResource()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Attributes())
	assert.Greater(t, res.Len(), 0)
}

func TestResourceServiceNameRequired(t *testing.T) {
	os.Clearenv()
	initConfig()
	res, err := newResource()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Zero(t, res)
}

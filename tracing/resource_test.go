package tracing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	svcname, svcver := "test-service", "1.0.0"
	res, _ := newResource(svcname, svcver)

	assert.NotNil(t, res)
	assert.NotNil(t, res.Attributes())
	assert.Greater(t, res.Len(), 0)
}

func TestResourceServiceNameRequired(t *testing.T) {
	svcname, svcver := "", "1.0.0"
	res, err := newResource(svcname, svcver)

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Zero(t, res)
}

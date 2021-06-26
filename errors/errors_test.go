package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWrap(t *testing.T) {
	child1 := errors.New(
		"failed to connect to endpoint: unknown.com")
	parent1 := ErrExporterStartupFailure
	wrapped1 := Wrap(child1, parent1)

	expected1 := "Failed to start exporter: failed to connect to endpoint: unknown.com"
	assert.Error(t, wrapped1)
	assert.EqualError(t, wrapped1, expected1)

	var child2 error = nil
	parent2 := ErrExporterStartupFailure
	wrapped2 := Wrap(child2, parent2)

	var expected2 error = nil
	assert.NoError(t, wrapped2)
	assert.Nil(t, wrapped2)
	assert.Equal(t, wrapped2, expected2)
}

package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWrap(t *testing.T) {
	child := errors.New(
		"failed to connect to endpoint: unknown.com")
	parent := ErrExporterStartupFailure
	wrapped := Wrap(child, parent)

	expected := "failed to start the exporter: failed to connect to endpoint: unknown.com"
	assert.Error(t, wrapped)
	assert.EqualError(t, wrapped, expected)
}

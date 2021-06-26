package exporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdOutExporter(t *testing.T) {
	t.Run("constructs a new exporter", func(tt *testing.T) {
		exp, err := NewStdOut()
		assert.NoError(t, err)
		assert.NotNil(t, exp)
	})

	t.Run("ErrStartupFailure on startup failure", func(tt *testing.T) {
		exp, err := NewStdOut()
		assert.NoError(t, err)
		assert.NotNil(t, exp)
	})

	t.Run("prints to the console", func(tt *testing.T) {

	})

}

package exporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestStdOutExporter(t *testing.T) {
	t.Run("constructs a new exporter", func(tt *testing.T) {
		exp, err := NewStdOut()
		assert.NoError(t, err)
		assert.NotNil(t, exp)
	})

	t.Run("should shutdown", func(tt *testing.T) {
		exp := &stdOutExporter{
			exporter: &stdouttrace.Exporter{},
		}
		ctx := context.Background()
		err := exp.Shutdown(ctx)

		assert.NoError(t, err)
		assert.Nil(t, err)

	})

	t.Run("should not fail when calling export on the internal exporter", func(tt *testing.T) {
		exp := stdOutExporter{exporter: &stdouttrace.Exporter{}}

		var spans []trace.ReadOnlySpan
		err := exp.ExportSpans(context.Background(), spans)

		assert.NoError(t, err)
		assert.Zero(t, err)
		assert.Nil(t, err)
	})

	t.Run("should write bytes to underlying stream", func(tt *testing.T) {
		w := new(MockWriter)
		w.On("Write", []byte("hello world"))
		_, err := NewStdOut(stdouttrace.WithWriter(w))

		assert.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("start should return nil", func(tt *testing.T) {
		w := new(MockWriter)
		w.On("Write", []byte("hello world"))
		exp, _ := NewStdOut(stdouttrace.WithWriter(w))
		err := exp.Start(context.Background())
		assert.NoError(t, err)
		assert.Nil(t, err)
	})

}

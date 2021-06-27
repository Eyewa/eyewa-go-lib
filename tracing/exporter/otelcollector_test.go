package exporter

import (
	"context"
	"errors"
	"testing"

	liberrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/stretchr/testify/assert"
)

func TestOtelCollectorUsingOtelExporter(t *testing.T) {
	t.Run("should error on start failure", func(tt *testing.T) {
		mock := new(OtelColMock)
		ctx := context.Background()
		mock.On("Start", ctx).Return(errors.New("failed to start"))

		ep := "test.endpoint.com"
		exp := &otelCollectorExporter{exporter: mock, endpoint: ep}
		err := exp.Start(ctx)

		expected := liberrs.ErrExporterStartupFailure.Error() + ":" + err.Error()
		assert.EqualError(t, err, expected)
		assert.Error(t, err)
		assert.NotZero(t, err)
	})
}

// func TestNewCollector(t *testing.T) {
// 	ep := "cool.endpoint.com"
// 	exp, err := NewOpenTelemetryCollectorExporter(ep, true, false)
// 	fmt.Println(err)
// }

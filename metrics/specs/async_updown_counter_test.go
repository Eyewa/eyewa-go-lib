package specs

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/expfmt"
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

var _ = Describe("Given that metric launcher is launched", func() {
	var (
		expectedInstrumentName      = "test_async_updown_counter"
		expectedInstrumentVersion   = "1.0.0"
		expectedValue               = 35.0
		expectedInstrumentationType = "GAUGE"
	)

	Describe(fmt.Sprintf("When async updown counter is initialized and being increased with value %f", expectedValue), func() {
		It("should return expected metric result", func() {
			callback := func(ctx context.Context, result metric.Float64ObserverResult) {
				result.Observe(expectedValue)
			}

			_ = meter.NewAsyncUpDownCounter(expectedInstrumentName,
				callback,
				metric.WithInstrumentationVersion(expectedInstrumentVersion),
			)

			res, err := http.Get(URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			if v, ok := mf[expectedInstrumentName]; ok {
				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				actualValue := v.GetMetric()[0].Gauge.Value
				Expect(*actualValue).Should(Equal(expectedValue))
			} else {
				Fail("Measurement couldn't find")
			}
		})
	})
})

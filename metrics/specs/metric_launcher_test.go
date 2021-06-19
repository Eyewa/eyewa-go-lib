package specs

import (
	"fmt"
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/expfmt"
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

var _ = Describe("Given that metric launcher is launched", func() {
	var (
		expectedCounterName         = "test_counter"
		expectedInstrumentVersion   = "1.0.0"
		expectedValue               = 45.0
		expectedInstrumentationType = "COUNTER"
	)

	Describe("When request sent to metric server without any instrumentation", func() {
		var (
			res *http.Response
			err error
		)

		It("should return http status ok", func() {
			res, err = http.Get(ts.URL)

			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))
		})
	})

	meter := metrics.NewMeter("test.meter", nil)

	Describe(fmt.Sprintf("When counter is initialized and being increased with expected value %f", expectedValue), func() {
		It("should return counted metric result", func() {
			counter := meter.NewCounter(expectedCounterName,
				metric.WithInstrumentationVersion(expectedInstrumentVersion),
			)
			counter.Add(expectedValue)

			res, err := http.Get(ts.URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			for actualCounterName, v := range mf {
				Expect(actualCounterName).Should(Equal(expectedCounterName))

				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				actualValue := v.GetMetric()[0].Counter.Value
				Expect(actualValue).Should(Equal(expectedValue))
			}
		})
	})
})

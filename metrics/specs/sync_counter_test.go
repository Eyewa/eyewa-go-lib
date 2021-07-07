package specs

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/expfmt"
	"go.opentelemetry.io/otel/metric"
)

var _ = Describe("Given that metric launcher is launched", func() {
	var (
		expectedInstrumentName      = "test_counter"
		expectedInstrumentVersion   = "1.0.0"
		expectedValue               = 45.0
		expectedInstrumentationType = "COUNTER"
	)

	Describe(fmt.Sprintf("When counter is initialized and being increased with value %f", expectedValue), func() {
		It("should return expected metric result", func() {
			counter, err := meter.NewCounter(expectedInstrumentName,
				metric.WithInstrumentationVersion(expectedInstrumentVersion),
			)
			Expect(err).Should(BeNil())

			counter.Add(expectedValue)

			res, err := http.Get(URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			if v, ok := mf[expectedInstrumentName]; ok {
				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				actualValue := v.GetMetric()[0].Counter.Value
				Expect(*actualValue).Should(Equal(expectedValue))
			} else {
				Fail("Measurement couldn't find")
			}
		})
	})
})

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
		expectedInstrumentName      = "test_counter"
		expectedInstrumentVersion   = "1.0.0"
		expectedValue               = 45.0
		expectedInstrumentationType = "COUNTER"
	)

	meter := metrics.NewMeter("test.meter", nil)

	Describe(fmt.Sprintf("When counter is initialized and being increased with value %f", expectedValue), func() {
		It("should return expected metric result", func() {
			counter := meter.NewCounter(expectedInstrumentName,
				metric.WithInstrumentationVersion(expectedInstrumentVersion),
			)
			counter.Add(expectedValue)

			res, err := http.Get(ts.URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			isMeasurementExist := false
			for actualCounterName, v := range mf {
				fmt.Println(v)
				if actualCounterName != expectedInstrumentName {
					continue
				}

				isMeasurementExist = true
				Expect(actualCounterName).Should(Equal(expectedInstrumentName))

				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				actualValue := v.GetMetric()[0].Counter.Value
				Expect(*actualValue).Should(Equal(expectedValue))
			}

			Expect(isMeasurementExist).Should(Equal(true))
		})
	})
})

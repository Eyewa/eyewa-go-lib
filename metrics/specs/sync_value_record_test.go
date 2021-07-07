package specs

import (
	"net/http"
	"sync"

	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/expfmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var _ = Describe("Given that metric launcher is launched", func() {
	var (
		expectedInstrumentName      = "test_value_recorder"
		expectedInstrumentVersion   = "1.0.0"
		expectedInstrumentationType = "HISTOGRAM"
		valueRecorder               *metrics.ValueRecorder
		once                        sync.Once
		err                         error
	)

	BeforeEach(func() {
		once.Do(func() {
			valueRecorder, err = meter.NewValueRecorder(expectedInstrumentName,
				metric.WithInstrumentationVersion(expectedInstrumentVersion),
			)
			Expect(err).Should(BeNil())
		})
	})

	Describe("When value recorder is initialized and being increased 2 times", func() {
		It("should return expected metric result", func() {
			var (
				expectedMetricCount = 1
				expectedSampleCount = uint64(2)
				firstRecordValue    = 35.0
				secondRecordValue   = 55.0
			)
			// First metric, first record
			valueRecorder.Record(firstRecordValue, attribute.Any("Name", "FirstMetric"))
			// First metric, second record
			valueRecorder.Record(secondRecordValue, attribute.Any("Name", "FirstMetric"))

			res, err := http.Get(URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			if v, ok := mf[expectedInstrumentName]; ok {
				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				metrics := v.GetMetric()
				Expect(metrics).Should(HaveLen(expectedMetricCount))

				metric := v.GetMetric()[0]

				actualSampleCount := metric.Histogram.SampleCount
				Expect(*actualSampleCount).Should(Equal(expectedSampleCount))

				actualSampleSum := metric.Histogram.SampleSum
				Expect(*actualSampleSum).Should(Equal(firstRecordValue + secondRecordValue))
			} else {
				Fail("Measurement couldn't find")
			}
		})
	})

	Describe("When value recorder being increased with different attribute", func() {
		It("should return another expected metric result", func() {
			var (
				expectedMetricCount = 2
				expectedSampleCount = uint64(1)
				firstRecordValue    = 40.0
			)

			// Second metric, first record
			valueRecorder.Record(firstRecordValue, attribute.Any("Name", "SecondMetric"))

			res, err := http.Get(URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			if v, ok := mf[expectedInstrumentName]; ok {
				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				metrics := v.GetMetric()
				Expect(metrics).Should(HaveLen(expectedMetricCount))

				// pick second metric
				metric := v.GetMetric()[1]

				actualSampleCount := metric.Histogram.SampleCount
				Expect(*actualSampleCount).Should(Equal(expectedSampleCount))

				actualSampleSum := metric.Histogram.SampleSum
				Expect(*actualSampleSum).Should(Equal(firstRecordValue))
			} else {
				Fail("Measurement couldn't find")
			}
		})
	})
})

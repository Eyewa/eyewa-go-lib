package specs

import (
	"fmt"
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/expfmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"sync"
)

var _ = Describe("Given that metric launcher is launched", func() {
	var (
		expectedInstrumentName      = "test_value_recorder"
		expectedInstrumentVersion   = "1.0.0"
		firstRecordValue            = 35.0
		secondRecordValue           = 55.0
		expectedInstrumentationType = "HISTOGRAM"
		valueRecorder               metrics.ValueRecorder
		once                        sync.Once
	)

	BeforeEach(func() {
		once.Do(func() {
			valueRecorder = meter.NewValueRecorder(expectedInstrumentName,
				metric.WithInstrumentationVersion(expectedInstrumentVersion),
			)
		})
	})

	Describe("When value recorder is initialized and being increased 2 times", func() {
		It("should return expected metric result", func() {
			var (
				expectedMetricCount = 1
				expectedSampleCount = uint64(2)
			)
			// first record
			valueRecorder.Record(firstRecordValue, attribute.Any("Name", "FirstMetric"))
			// second record
			valueRecorder.Record(secondRecordValue, attribute.Any("Name", "FirstMetric"))

			res, err := http.Get(ts.URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			if v, ok := mf[expectedInstrumentName]; ok {
				fmt.Println(v)
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
			)

			valueRecorder.Record(firstRecordValue, attribute.Any("Name", "SecondMetric"))

			res, err := http.Get(ts.URL)
			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(res.Body)
			Expect(err).Should(BeNil())

			if v, ok := mf[expectedInstrumentName]; ok {
				fmt.Println(v)
				actualInstrumentationType := v.GetType().String()
				Expect(actualInstrumentationType).Should(Equal(expectedInstrumentationType))

				metrics := v.GetMetric()
				Expect(metrics).Should(HaveLen(expectedMetricCount))

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

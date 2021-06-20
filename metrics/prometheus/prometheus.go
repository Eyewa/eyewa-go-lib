package prometheus

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"net/http"
	"time"
)

// ExportOption is configuration for MetricLauncher.
type ExportOption struct {
	// CollectPeriod sets period interval exporter.
	CollectPeriod time.Duration
}

// PrometheusExporter is a pull-based exporter
type PrometheusExporter struct {
	exportOption ExportOption
	exporter     *prometheus.Exporter
}

// NewPrometheusExporter creates new PrometheusExporter with given ExportOption
func NewPrometheusExporter(option ExportOption) (*PrometheusExporter, error) {
	config := prometheus.Config{}
	c := controller.New(
		processor.New(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(option.CollectPeriod),
	)

	exporter, err := prometheus.New(config, c)

	if err != nil {
		return nil, metrics.FailedToInitPrometheusExporterError.Inner(err)
	}

	return &PrometheusExporter{
		exportOption: option,
		exporter:     exporter,
	}, nil
}

// ServeHTTP implements http.Handler.
func (p *PrometheusExporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.exporter.ServeHTTP(w, r)
}

// MeterProvider returns the MeterProvider of this exporter.
func (p *PrometheusExporter) MeterProvider() metric.MeterProvider {
	return p.exporter.MeterProvider()
}

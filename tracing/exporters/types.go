package exporters

import (
	"go.opentelemetry.io/otel/exporters/otlp"
	stdout "go.opentelemetry.io/otel/exporters/stdout"
)

type stdOutExporter struct {
	exporter *stdout.Exporter
}

type otelCollectorExporter struct {
	exporter *otlp.Exporter
}

type otelExporterInterface struct {
}

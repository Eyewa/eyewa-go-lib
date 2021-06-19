package metrics

import (
	"fmt"
)

const MetricPrefix = "MET"

// Error is wrapper for error interface
type Error struct {
	prefix string
	id     int
	desc   string
	// Err is the underlying error.
	Err error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s%.2d - %s - Error is %v", e.prefix, e.id, e.desc, e.Err)
	} else {
		return fmt.Sprintf("%s%.2d - %s", e.prefix, e.id, e.desc)
	}
}

// Inner passes err as the underlying error
func (e *Error) Inner(err error) error {
	e.Err = err
	return e
}

func defineError(prefix string, id int, desc string) Error {
	return Error{
		prefix,
		id,
		desc,
		nil,
	}
}

var (
	PrometheusExporterInitFailedError = defineError(MetricPrefix, 1, "Failed to initialize prometheus exporter.")
	FailedToStartRuntimeMetricsError  = defineError(MetricPrefix, 2, "Failed to start runtime metrics.")
	FailedToStartHostMetricsError     = defineError(MetricPrefix, 3, "Failed to start host metrics.")
	FailedToStartMetricServerError    = defineError(MetricPrefix, 4, "Failed to start metric server.")
)

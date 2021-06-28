package errors

import (
	"errors"

	pkgErrs "github.com/pkg/errors"
)

var (
	// MessageBrokerClient errors
	ErrorNoQueuesSpecified           = errors.New("No queues specified! Cannot consume/publish to any queue.")
	ErrorNoConsumerQueueSpecified    = errors.New("No queue specified to consume from!")
	ErrorNoPublisherQueueSpecified   = errors.New("No queue specified to publish to!")
	ErrorNoRMQConnection             = errors.New("No connection to RMQ exists!")
	ErrorBrokerClientNotRecognized   = errors.New("Broker client not recognized.")
	ErrorFailedToPublishToDeadletter = errors.New("Failed to publish event error to deadletter queue.")
	ErrorFailedToPublishEvent        = errors.New("Failed to publish event to queue.")

	// Tracing errors
	ErrorNoExporterEndpointSpecified = errors.New("No exporter endpoint specified.")
	ErrorNoServiceNameSpecified      = errors.New("No service name specified.")

	// Metrics errors
	ErrorFailedToInitPrometheusExporter = errors.New("Failed to initialize prometheus exporter.")
	ErrorFailedToStartRuntimeMetrics    = errors.New("Failed to start runtime metrics.")
	ErrorFailedToStartHostMetrics       = errors.New("Failed to start host metrics.")
	ErrorFailedToStartMetricServer      = errors.New("Failed to start metric server.")
	ErrorFailedToCreateInstrument       = errors.New("Failed to create instrument.")

	// DBClient errors
	ErrorNoDBDriverSpecified          = errors.New("No DB driver specified.")
	ErrorUnsupportedDBDriverSpecified = errors.New("Unsupported DB driver specified.")
	ErrorNoDBClientFound              = errors.New("Failed to close connection. No db client found")
)

// Wrap wraps a child error with a parent.
func Wrap(child, parent error) error {
	if child == nil {
		return nil
	}
	return pkgErrs.Wrap(child, parent.Error())
}

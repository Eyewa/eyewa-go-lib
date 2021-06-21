package errors

import "errors"

var (
	ErrorNoQueuesSpecified           = errors.New("No queues specified! Cannot consume/publish to any queue.")
	ErrorNoConsumerQueueSpecified    = errors.New("No queue specified to consume from!")
	ErrorNoPublisherQueueSpecified   = errors.New("No queue specified to publish to!")
	ErrorNoRMQConnection             = errors.New("No connection to RMQ exists!")
	ErrorBrokerClientNoRecognized    = errors.New("Broker client not recognized.")
	ErrorFailedToPublishToDeadletter = errors.New("Failed to publish event error to deadletter queue.")
	ErrorFailedToPublishEvent        = errors.New("Failed to publish event to queue.")

	FailedToInitPrometheusExporterError = errors.New("Failed to initialize prometheus exporter.")
	FailedToStartRuntimeMetricsError    = errors.New("Failed to start runtime metrics.")
	FailedToStartHostMetricsError       = errors.New("Failed to start host metrics.")
	FailedToStartMetricServerError      = errors.New("Failed to start metric server.")
)

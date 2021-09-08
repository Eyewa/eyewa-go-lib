package errors

import (
	"errors"
)

var (
	// MessageBrokerClient errors
	ErrorNoQueuesSpecified               = errors.New("No queues specified! Cannot consume/publish to any queue.")
	ErrorNoConsumerQueueSpecified        = errors.New("No queue specified to consume from!")
	ErrorNoPublisherQueueSpecified       = errors.New("No queue specified to publish to!")
	ErrorNoRMQConnection                 = errors.New("No connection to RMQ exists!")
	ErrorChannelDoesNotExist             = errors.New("Channel does not exist!")
	ErrorBrokerClientNotRecognized       = errors.New("Broker client not recognized.")
	ErrorFailedToPublishToDeadletter     = errors.New("Failed to publish event error to deadletter queue.")
	ErrorFailedToPublishEvent            = errors.New("Failed to publish event to queue.")
	ErrorLostConnectionToMessageBroker   = errors.New("Lost connection to Message Broker!")
	ErrorAckFailure                      = errors.New("Failed to acknowledge new message delivered to client")
	ErrorNackFailure                     = errors.New("Failed to unacknowledge message.")
	ErrorConsumeFailure                  = errors.New("Failed to consume from queue(%s). %s")
	ErrorEventUnmarshalFailure           = errors.New("Failed to unmarshal event from queue(%s). %s")
	ErrorQueueDeclareFailure             = errors.New("Failed to declare queue(%s). %s")
	ErrorExchangeDeclareFailure          = errors.New("Failed to declare an exchange for queue(%s). %s")
	ErrorExchangeBindFailure             = errors.New("Failed to bind exchange to queue(%s). %s")
	ErrorChannelCreateFailure            = errors.New("Failed to create new channel for queue(%s). %s")
	ErrorQueueInspectFailure             = errors.New("Failed to inspect queue(%s). %s")
	ErrorQueueInspectMissingQueueFailure = errors.New("Queue specified to inspect doesn't exist queue(%s)")

	// Tracing errors
	ErrorNoExporterEndpointSpecified = errors.New("No exporter endpoint specified.")
	ErrorNoServiceNameSpecified      = errors.New("No service name specified.")

	// Metrics errors
	ErrorFailedToInitPrometheusExporter = errors.New("Failed to initialize prometheus exporter.")
	ErrorFailedToStartRuntimeMetrics    = errors.New("Failed to start runtime metrics.")
	ErrorFailedToStartHostMetrics       = errors.New("Failed to start host metrics.")
	ErrorFailedToStartMetricServer      = errors.New("Failed to start metric server. Error is %s")
	ErrorFailedToCreateInstrument       = errors.New("Failed to create instrument.")

	// DBClient errors
	ErrorNoDBDriverSpecified          = errors.New("No DB driver specified.")
	ErrorUnsupportedDBDriverSpecified = errors.New("Unsupported DB driver specified.")
	ErrorNoDBClientFound              = errors.New("Failed to close connection. No db client found.")
)

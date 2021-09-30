package rabbitmq

import (
	"sync"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

// Config for all RabbitMQ env vars
type Config struct {
	MessageBroker string `mapstructure:"message_broker"`

	// RMQ credentials
	Server   string `mapstructure:"rabbitmq_server"`
	AmqpPort string `mapstructure:"rabbitmq_amqp_port"`
	Username string `mapstructure:"rabbitmq_username"`
	Password string `mapstructure:"rabbitmq_password"`

	// No. of messsages RMQ should send to a consumer
	// https://www.cloudamqp.com/blog/how-to-optimize-the-rabbitmq-prefetch-count.html
	QueuePrefetchCount string `mapstructure:"queue_prefetch_count"`

	// Queues for consuming + publishing
	PublisherQueueName string `mapstructure:"publisher_queue_name"`
	ConsumerQueueName  string `mapstructure:"consumer_queue_name"`

	ConsumerExchange string `mapstructure:"rabbitmq_consumer_exchange"`

	// Exchanges to bind consumer + publisher queues to
	PublisherExchangeType string `mapstructure:"rabbitmq_publisher_exchange_type"`
	ConsumerExchangeType  string `mapstructure:"rabbitmq_consumer_exchange_type"`

	// Purely for identifying what service/service instance is connected to a RMQ channel
	ServiceName string `mapstructure:"service_name"`
	HostName    string
}

// RMQClient RMQ client for implementing the MessageBroker interface and handling all things RMQ.
type RMQClient struct {
	mutex *sync.RWMutex

	connection *amqp.Connection

	// Map of channels for all queues
	channels map[string]*amqp.Channel
}

type unmarshalledEyewaEvent struct {
	unmarshalledCommon
	event    *base.EyewaEvent
	callback base.MessageBrokerCallbackFunc
}

type unmarshalledMagentoEvent struct {
	unmarshalledCommon
	event    *base.MagentoProductEvent
	callback base.MessageBrokerMagentoProductCallbackFunc
}

type unmarshalledCommon struct {
	queue   string
	msg     amqp.Delivery
	span    trace.Span
	started time.Time
	err     error
}

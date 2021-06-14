package rabbitmq

import (
	"sync"

	"github.com/streadway/amqp"
)

// Config for all RabbitMQ env vars
type Config struct {
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

	// Exchanges to bind consumer + publisher queues to
	PublisherExchangeType string `mapstructure:"rabbitmq_publisher_exchange_type"`
	ConsumerExchangeType  string `mapstructure:"rabbitmq_consumer_exchange_type"`
}

// RMQClient RMQ client for implementing the MessageBroker interface and handling all things RMQ.
type RMQClient struct {
	mutex *sync.RWMutex

	connection *amqp.Connection

	// map of channels for all queues
	channels map[string]*amqp.Channel
}

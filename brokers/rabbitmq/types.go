package rabbitmq

import "github.com/streadway/amqp"

// Config for all RabbitMQ env vars
type Config struct {
	Server   string `mapstructure:"rabbitmq_server"`
	Port     string `mapstructure:"rabbitmq_amqp_port"`
	Username string `mapstructure:"rabbitmq_username"`
	Password string `mapstructure:"rabbitmq_password"`
}

// RMQClient RMQ client for implementing the MessageBroker interface and handling all things RMQ.
type RMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

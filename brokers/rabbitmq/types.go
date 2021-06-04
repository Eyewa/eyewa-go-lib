package rabbitmq

// EnvConfig for all RabbitMQ env vars
type EnvConfig struct {
	RabbitMQServer    string `env:"RABBITMQ_SERVER,required"`
	RabbitMQPort      int    `env:"RABBITMQ_AMQP_PORT,required"`
	RabbitMQAdminPort int    `env:"RABBITMQ_ADMIN_PORT,required"`
	RabbitMQUsername  string `env:"RABBITMQ_USERNAME,required"`
	RabbitMQPassword  string `env:"RABBITMQ_PASSWORD,required"`
}

// RMQClient RMQ client for implementing the MessageBroker interface and handling all things RMQ.
type RMQClient struct{}

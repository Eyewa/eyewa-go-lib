package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/ory/viper"
	"github.com/streadway/amqp"
)

var config Config

// NewRMQClient new rmq client
func NewRMQClient() *RMQClient {
	return new(RMQClient)
}

// Connect Establishes connnection to the message broker of choice
func (rmq *RMQClient) Connect() error {
	_, connStr, err := initConfig()
	if err != nil {
		return err
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return err
	}

	rmq.connection = conn

	rmq.channel, err = conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

// Publish publishes a message to a queue
func (rmq *RMQClient) Publish(queue string) error {
	return nil
}

// CloseConnection closes a connection as well as any
// underlying channels associated to it.
func (rmq *RMQClient) CloseConnection() error {
	if rmq.connection != nil {
		return rmq.connection.Close()
	}

	return nil
}

// Consume consumes messages from a queue
func (rmq *RMQClient) Consume(queue string) error {
	return nil
}

func initConfig() (Config, string, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	envVars := []string{
		"RABBITMQ_SERVER",
		"RABBITMQ_AMQP_PORT",
		"RABBITMQ_USERNAME",
		"RABBITMQ_PASSWORD",
	}

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return config, "", err
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, "", err
	}

	return config, fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password,
		config.Server, config.Port), nil
}

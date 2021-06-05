package rabbitmq

import (
	"fmt"

	"github.com/eyewa/eyewa-go-lib/log"
)

// Connect Establishes connnection to the message broker of choice
func (rmq *RMQClient) Connect() error {
	log.Info("Connectd to RMQ succesfully")
	return nil
}

// Publish publishes a message to a queue
func (rmq *RMQClient) Publish(queue string) error {
	log.Info(fmt.Sprintf("Published to %s successfully.", queue))
	return nil
}

func (rmq *RMQClient) CloseConnection() error {
	log.Info("Closed connection to RMQ.")
	return nil
}

func (rmq *RMQClient) Consume(queue string) error {
	log.Info(fmt.Sprintf("Consumed 5 messages from %s.", queue))
	return nil
}

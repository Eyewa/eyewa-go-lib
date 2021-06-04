package brokers

import (
	"errors"
	"strings"

	"github.com/eyewa/eyewa-go-lib/brokers/kafka"
	"github.com/eyewa/eyewa-go-lib/brokers/rabbitmq"
	"github.com/eyewa/eyewa-go-lib/brokers/sqs"
)

const (
	RabbitMQ BrokerType = "rabbitmq"
	SQS      BrokerType = "sqs"
	Kafka    BrokerType = "kafka"
)

var (
	broker *MessageBrokerClient
)

// OpenConnection opens a connection to the message broker
func OpenConnection(brokerType string) error {
	switch strings.ToLower(brokerType) {
	case string(RabbitMQ):
		broker = &MessageBrokerClient{RabbitMQ, new(rabbitmq.RMQClient)}
	case string(SQS):
		broker = &MessageBrokerClient{SQS, new(sqs.SQSClient)}
	case string(Kafka):
		broker = &MessageBrokerClient{Kafka, new(kafka.KafkaClient)}
	}

	return connect(broker)
}

func connect(broker *MessageBrokerClient) error {
	if broker.Client != nil {
		err := broker.Client.Connect()
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("Client not recognized.")
}

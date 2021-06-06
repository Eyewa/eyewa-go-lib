package brokers

import (
	"errors"
	"os"
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
func OpenConnection() (*MessageBrokerClient, error) {
	switch strings.ToLower(os.Getenv("MESSAGE_BROKER")) {
	case string(RabbitMQ):
		broker = &MessageBrokerClient{RabbitMQ, new(rabbitmq.RMQClient)}
	case string(SQS):
		broker = &MessageBrokerClient{SQS, new(sqs.SQSClient)}
	case string(Kafka):
		broker = &MessageBrokerClient{Kafka, new(kafka.KafkaClient)}
	default:
		broker = new(MessageBrokerClient)
	}

	return broker.connect()
}

func (*MessageBrokerClient) connect() (*MessageBrokerClient, error) {
	if broker.Client != nil {
		err := broker.Client.Connect()
		if err != nil {
			return nil, err
		}

		return broker, nil
	}

	return nil, errors.New("Client not recognized.")
}

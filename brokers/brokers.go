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

func NewMessageBrokerClient(brokerType BrokerType, client MessageBroker) *MessageBrokerClient {
	return &MessageBrokerClient{brokerType, client}
}

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

// NewConsumerClient creates a new consumer client
func NewConsumerClient(brokerType BrokerType) *MessageBrokerConsumerClient {
	client := getClient(brokerType)
	if client != nil {
		return &MessageBrokerConsumerClient{brokerType, client}
	}

	return new(MessageBrokerConsumerClient)
}

// NewPublisherClient creates a new publisher client
func NewPublisherClient(brokerType BrokerType) *MessageBrokerPublisherClient {
	client := getClient(brokerType)
	if client != nil {
		return &MessageBrokerPublisherClient{brokerType, client}
	}

	return new(MessageBrokerPublisherClient)
}

func getClient(brokerType BrokerType) MessageBroker {
	clientMap := map[BrokerType]MessageBroker{
		RabbitMQ: rabbitmq.NewRMQClient(),
		SQS:      sqs.NewSQSClient(),
		Kafka:    kafka.NewKafkaClient(),
		Mock:     NewMockClient(),
	}

	if broker, ok := clientMap[brokerType]; ok {
		return broker
	}

	return nil
}

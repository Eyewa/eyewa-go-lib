package brokers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/eyewa/eyewa-go-lib/brokers/kafka"
	"github.com/eyewa/eyewa-go-lib/brokers/rabbitmq"
	"github.com/eyewa/eyewa-go-lib/brokers/sqs"
	libErrs "github.com/eyewa/eyewa-go-lib/errors"
)

const (
	RabbitMQ BrokerType = "rabbitmq"
	SQS      BrokerType = "sqs"
	Kafka    BrokerType = "kafka"
)

var (
	broker               *MessageBrokerClient
	maxConnectionRetries uint64 = 100
)

func NewMessageBrokerClient(brokerType BrokerType, client MessageBroker) *MessageBrokerClient {
	return &MessageBrokerClient{brokerType, client, maxConnectionRetries}
}

// OpenConnection opens a connection to the message broker
func OpenConnection() (*MessageBrokerClient, error) {
	switch strings.ToLower(os.Getenv("MESSAGE_BROKER")) {
	case string(RabbitMQ):
		broker = &MessageBrokerClient{RabbitMQ, new(rabbitmq.RMQClient), maxConnectionRetries}
	case string(SQS):
		broker = &MessageBrokerClient{SQS, new(sqs.SQSClient), maxConnectionRetries}
	case string(Kafka):
		broker = &MessageBrokerClient{Kafka, new(kafka.KafkaClient), maxConnectionRetries}
	default:
		broker = new(MessageBrokerClient)
	}

	return broker.connect()
}

func (*MessageBrokerClient) connect() (*MessageBrokerClient, error) {
	if broker.Client != nil {
		// apply exponential backoff to try establishing a connection.
		connect := func() error {
			return broker.Client.Connect()
		}

		bkoff := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), broker.MaxConnectionRetries)
		err := backoff.RetryNotify(connect, bkoff, func(err error, duration time.Duration) {
			if err != nil {
				fmt.Println(err.Error())
			}
		})
		if err != nil {
			return nil, err
		}

		return broker, err
	}

	return nil, libErrs.ErrorBrokerClientNotRecognized
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

// AlwaysReconnect attempts to re-establish connection to
// a message broker using an exponential backoff.
//
// Consuming from a message broker should be a long lived connection
// Should it be lost for whatever reason, this func initiates the
// attempt of re-gaining it so consumption can resume.
//
// If for any reason connection is lost, the errCh is used to communicate
// with the callback on such an event. This ALSO has to be a goroutine
// so it is non-blocking as this is merely a retry strategy.
func AlwaysReconnect(errCh chan error, callback ConsumerCallbackFunc) {
	go func() {
		for {
			for err := range errCh {
				if err != nil {
					broker, err = OpenConnection()
					if err == nil {
						callback(broker, errCh)
					}
				}
			}
		}
	}()
}

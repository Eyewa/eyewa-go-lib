package brokers

import (
	"context"
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
)

// BrokerType represents a type of broker - sqs, rmq etc.
type BrokerType string

type ConsumerCallbackFunc func(broker *MessageBrokerClient, errCh chan error)

// MessageBrokerClient a message broker client with the
// capability to act as both consumer and publisher.
type MessageBrokerClient struct {
	Type                 BrokerType
	Client               MessageBroker
	MaxConnectionRetries uint64
}

// MessageBrokerConsumerClient a consumer client.
type MessageBrokerConsumerClient struct {
	Type   BrokerType
	Client Consumer
}

// MessageBrokerPublisherClient a publisher client.
type MessageBrokerPublisherClient struct {
	Type   BrokerType
	Client Publisher
}

// MessageBroker a contract any broker client should fulfil.
// Any client implementing this contract is assumed to be both
// a publisher and a consumer.
type MessageBroker interface {
	Connect() error
	IsConnectionOpen() bool
	CloseConnection() error
	ConnectionListener()
	Consumer
	Publisher
}

// Consumer a contract any consumer should fulfil.
type Consumer interface {
	Connect() error
	CloseConnection() error
	Consume(queue string, callback base.MessageBrokerCallbackFunc)
	ConsumeMagentoCatalog(queue string, callback base.MessageBrokerMagentoCatalogCallbackFunc)
}

// Publisher a contract any publisher should fulfil.
type Publisher interface {
	Connect() error
	CloseConnection() error
	Publish(ctx context.Context, queue string, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup)
}

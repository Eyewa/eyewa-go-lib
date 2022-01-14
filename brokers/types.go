package brokers

import (
	"context"
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/streadway/amqp"
)

// Queue priorities for declaring queues and publishing
const (
	PriorityNone = iota
	PriorityLow
	PriorityMedium
	PriorityMediumHigh
	PriorityHigh
	PriorityCritical
)

// ConnectFunc is the function that starts consuming from the given broker.
type ConsumeFunc func(broker *MessageBrokerClient) error

type ConsumerCallbackFunc func(broker *MessageBrokerClient, errCh chan error)

// BrokerType represents a type of broker - sqs, rmq etc.
type BrokerType string

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
	SendToDeadletterQueue(amqp.Delivery, error) error
	ConsumeMagentoProductEvents(queue string, callback base.MessageBrokerMagentoProductCallbackFunc)
}

// Publisher a contract any publisher should fulfil.
type Publisher interface {
	Connect() error
	CloseConnection() error
	Publish(ctx context.Context, queue string, priority int, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup)
	PublishMagentoProductEvent(ctx context.Context, queue string, priority int, event *base.MagentoProductEvent, callback base.MessageBrokerMagentoProductCallbackFunc, wg *sync.WaitGroup)
	PublishEvent(ctx context.Context, queue string, priority int, event *[]byte, wg *sync.WaitGroup) error
}

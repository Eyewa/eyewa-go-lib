package brokers

type BrokerType string

type MessageBrokerClient struct {
	Type   BrokerType
	Client MessageBroker
}

type MessageBrokerConsumerClient struct {
	Client Consumer
}

type MessageBrokerPublisherClient struct {
	Client Publisher
}

type MessageBroker interface {
	Connect() error
	CloseConnection() error
	Consumer
	Publisher
}

type Consumer interface {
	Connect() error
	Consume(queue string) error
	CloseConnection() error
}

type Publisher interface {
	Connect() error
	Publish(queue string) error
	CloseConnection() error
}

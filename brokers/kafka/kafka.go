package kafka

import (
	"context"
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
)

// NewKafkaClient new rmq client
func NewKafkaClient() *KafkaClient {
	return new(KafkaClient)
}

func (kafka *KafkaClient) Connect() error {
	return nil
}

func (kafka *KafkaClient) ConnectionListener() {
	//
}

func (kafka *KafkaClient) Publish(ctx context.Context, queue string, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup) {
	//
}

func (kafka *KafkaClient) CloseConnection() error {
	return nil
}

func (kafka *KafkaClient) Consume(queue string, callback base.MessageBrokerCallbackFunc) {
	//
}

func (kafka *KafkaClient) IsConnectionOpen() bool {
	return false
}

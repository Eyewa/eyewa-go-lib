package kafka

import (
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

func (kafka *KafkaClient) Publish(queue string, event *base.EyewaEvent, errChan chan<- error, wg *sync.WaitGroup) {
}

func (kafka *KafkaClient) CloseConnection() error {
	return nil
}

func (kafka *KafkaClient) Consume(queue string, callback base.ConsumeCallbackFunc) {
}

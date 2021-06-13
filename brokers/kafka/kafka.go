package kafka

import "sync"

// NewKafkaClient new rmq client
func NewKafkaClient() *KafkaClient {
	return new(KafkaClient)
}

func (kafka *KafkaClient) Connect() error {
	return nil
}

func (kafka *KafkaClient) Publish(queue string) error {
	return nil
}

func (kafka *KafkaClient) CloseConnection() error {
	return nil
}

func (kafka *KafkaClient) Consume(wg *sync.WaitGroup, queue string, errChan chan<- error) {
}

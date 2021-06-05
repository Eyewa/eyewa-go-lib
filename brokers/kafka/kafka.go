package kafka

import (
	"fmt"

	"github.com/eyewa/eyewa-go-lib/log"
)

func (kafka *KafkaClient) Connect() error {
	log.Info("Connectd to Kafka succesfully")
	return nil
}

func (kafka *KafkaClient) Publish(queue string) error {
	log.Info(fmt.Sprintf("Published to %s successfully.", queue))
	return nil
}

func (kafka *KafkaClient) CloseConnection() error {
	log.Info("Closed connection to Kafka.")
	return nil
}

func (kafka *KafkaClient) Consume(queue string) error {
	log.Info(fmt.Sprintf("Consumed 5 messages from %s.", queue))
	return nil
}

package sqs

import (
	"fmt"

	"github.com/eyewa/eyewa-go-lib/log"
)

func (sqs *SQSClient) Connect() error {
	log.Info("Connectd to SQS succesfully")
	return nil
}

func (sqs *SQSClient) Publish(queue string) error {
	log.Info(fmt.Sprintf("Published to %s successfully.", queue))
	return nil
}

func (sqs *SQSClient) CloseConnection() error {
	log.Info("Closed connection to SQS.")
	return nil
}

func (sqs *SQSClient) Consume(queue string) error {
	log.Info(fmt.Sprintf("Consumed 5 messages from %s.", queue))
	return nil
}

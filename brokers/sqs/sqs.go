package sqs

import (
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
)

// NewSQSClient new sqs client
func NewSQSClient() *SQSClient {
	return new(SQSClient)
}

func (sqs *SQSClient) Connect() error {
	return nil
}

func (sqs *SQSClient) Publish(queue string, event *base.EyewaEvent, errChan chan<- error, wg *sync.WaitGroup) {
}

func (sqs *SQSClient) CloseConnection() error {
	return nil
}

func (sqs *SQSClient) Consume(queue string, callback base.ConsumeCallbackFunc) {
}

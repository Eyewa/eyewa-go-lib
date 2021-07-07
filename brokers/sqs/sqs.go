package sqs

import (
	"context"
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

func (sqs *SQSClient) ConnectionListener() {
	//
}

func (sqs *SQSClient) Publish(ctx context.Context, queue string, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup) {
}

func (sqs *SQSClient) CloseConnection() error {
	return nil
}

func (sqs *SQSClient) Consume(queue string, callback base.MessageBrokerCallbackFunc) {
	//
}

func (sqs *SQSClient) IsConnectionOpen() bool {
	return false
}

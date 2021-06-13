package sqs

import "sync"

// NewSQSClient new sqs client
func NewSQSClient() *SQSClient {
	return new(SQSClient)
}

func (sqs *SQSClient) Connect() error {
	return nil
}

func (sqs *SQSClient) Publish(queue string) error {
	return nil
}

func (sqs *SQSClient) CloseConnection() error {
	return nil
}

func (sqs *SQSClient) Consume(wg *sync.WaitGroup, queue string, errChan chan<- error) {
}

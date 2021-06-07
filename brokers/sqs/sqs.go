package sqs

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

func (sqs *SQSClient) Consume(queue string) error {
	return nil
}

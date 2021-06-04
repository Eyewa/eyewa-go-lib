package rabbitmq

// Connect Establishes connnection to the message broker of choice
func (rmq *RMQClient) Connect() error {
	return nil
}

// Publish publishes a message to a queue
func (rmq *RMQClient) Publish(queue string) error {
	return nil
}

func (rmq *RMQClient) CloseConnection() error {
	return nil
}

func (rmq *RMQClient) Consume(queue string) error {
	return nil
}

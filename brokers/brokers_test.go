package brokers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	kafkaMock := new(ClientMock)
	kafkaMock.On("Connect").Return(nil)
	kafka := MessageBrokerClient{Kafka, kafkaMock}
	assert.NoError(t, connect(&kafka))

	rmqMock := new(ClientMock)
	rmqMock.On("Connect").Return(nil)
	rmq := MessageBrokerClient{RabbitMQ, rmqMock}
	assert.NoError(t, connect(&rmq))
}

func TestConnectFail(t *testing.T) {
	sqs := MessageBrokerClient{}
	sqs.Type = SQS
	assert.Error(t, connect(&sqs))

	rmqMock := new(ClientMock)
	rmqMock.On("Connect").Return(errors.New("bleh"))
	rmq := MessageBrokerClient{RabbitMQ, rmqMock}
	assert.Error(t, connect(&rmq))
}

func TestOpenConnection(t *testing.T) {
	rmqMock := new(ClientMock)
	rmqMock.On("Connect").Return(nil)
	_ = OpenConnection("rabbitmq")
	assert.Equal(t, RabbitMQ, broker.Type)

	sqsMock := new(ClientMock)
	sqsMock.On("Connect").Return(nil)
	_ = OpenConnection("sqs")
	assert.Equal(t, SQS, broker.Type)

	kafkaMock := new(ClientMock)
	kafkaMock.On("Connect").Return(nil)
	_ = OpenConnection("kafka")
	assert.Equal(t, Kafka, broker.Type)
}

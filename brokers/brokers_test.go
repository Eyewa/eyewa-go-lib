package brokers

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	var err error

	sqsMock := new(ClientMock)
	broker = &MessageBrokerClient{SQS, sqsMock}
	sqsMock.On("Connect").Return(nil)
	broker, err = broker.connect()
	assert.Nil(t, err)
	assert.NotNil(t, broker)
}

func TestConnectFail(t *testing.T) {
	var err error

	sqsMock := new(ClientMock)
	broker = &MessageBrokerClient{SQS, sqsMock}
	sqsMock.On("Connect").Return(errors.New("bleh"))
	broker, err = broker.connect()
	assert.EqualError(t, err, "bleh")
	assert.Nil(t, broker)
}

func TestOpenConnection(t *testing.T) {
	var err error

	os.Setenv("MESSAGE_BROKER", "rabbitmq")
	rmqMock := new(ClientMock)
	rmqMock.On("Connect").Return(nil)
	_, err = OpenConnection()
	assert.Equal(t, RabbitMQ, broker.Type)
	assert.NoError(t, err)

	os.Setenv("MESSAGE_BROKER", "sqs")
	sqsMock := new(ClientMock)
	sqsMock.On("Connect").Return(nil)
	_, err = OpenConnection()
	assert.Equal(t, SQS, broker.Type)
	assert.NoError(t, err)

	os.Setenv("MESSAGE_BROKER", "kafka")
	kafkaMock := new(ClientMock)
	kafkaMock.On("Connect").Return(nil)
	_, err = OpenConnection()
	assert.Equal(t, Kafka, broker.Type)
	assert.NoError(t, err)
}

func TestOpenConnectionFail(t *testing.T) {
	var err error

	os.Setenv("MESSAGE_BROKER", "activemq")
	rmqMock := new(ClientMock)
	rmqMock.AssertNotCalled(t, "Connect")
	_, err = OpenConnection()
	assert.Error(t, err)
	assert.Equal(t, "Client not recognized.", err.Error())

	os.Setenv("MESSAGE_BROKER", "")
	rmqMock = new(ClientMock)
	rmqMock.AssertNotCalled(t, "Connect")
	_, err = OpenConnection()
	assert.Error(t, err)
	assert.Equal(t, "Client not recognized.", err.Error())
}

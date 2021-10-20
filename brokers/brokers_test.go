package brokers

import (
	"errors"
	"os"
	"testing"

	"github.com/eyewa/eyewa-go-lib/brokers/rabbitmq"
	libErrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	var err error

	sqsMock := new(ClientMock)
	broker = &MessageBrokerClient{RabbitMQ, sqsMock, 1}
	sqsMock.On("Connect").Return(nil)
	broker, err = broker.connect()
	assert.Nil(t, err)
	assert.NotNil(t, broker)
}

func TestConnectFail(t *testing.T) {
	var err error

	sqsMock := new(ClientMock)
	broker = &MessageBrokerClient{RabbitMQ, sqsMock, 1}
	sqsMock.On("Connect").Return(errors.New("bleh"))
	broker, err = broker.connect()
	assert.EqualError(t, err, "bleh")
	assert.Nil(t, broker)
}

func TestOpenConnection(t *testing.T) {
	var err error

	rmqMock := new(ClientMock)
	rmqMock.On("Connect").Return(nil)
	_, err = OpenMockConnection(rmqMock)
	assert.Equal(t, Mock, broker.Type)
	assert.NoError(t, err)
}

func TestOpenConnectionFail(t *testing.T) {
	var err error

	os.Setenv("MESSAGE_BROKER", "activemq")
	rmqMock := new(ClientMock)
	rmqMock.AssertNotCalled(t, "Connect")
	_, err = OpenConnection()
	assert.Error(t, err)
	assert.EqualError(t, libErrs.ErrorBrokerClientNotRecognized, err.Error())

	os.Setenv("MESSAGE_BROKER", "")
	rmqMock = new(ClientMock)
	rmqMock.AssertNotCalled(t, "Connect")
	_, err = OpenConnection()
	assert.Error(t, err)
	assert.EqualError(t, libErrs.ErrorBrokerClientNotRecognized, err.Error())
}

func TestGetClient(t *testing.T) {
	var client MessageBroker

	client = getClient(RabbitMQ)
	assert.NotZero(t, client)

	client = getClient(Mock)
	assert.NotZero(t, client)
}

func TestNewConsumerClient(t *testing.T) {
	var client *MessageBrokerConsumerClient

	client = NewConsumerClient(RabbitMQ)
	assert.NotNil(t, client.Client)
	assert.Equal(t, RabbitMQ, client.Type)

	client = NewConsumerClient(Mock)
	assert.NotNil(t, client.Client)
	assert.Equal(t, Mock, client.Type)

	var activemq BrokerType = "activemq"
	client = NewConsumerClient(activemq)
	assert.Nil(t, client.Client)
	assert.Empty(t, client.Type)
}

func TestNewPublisherClient(t *testing.T) {
	var pub *MessageBrokerPublisherClient

	pub = NewPublisherClient(RabbitMQ)
	assert.NotNil(t, pub.Client)
	assert.Equal(t, RabbitMQ, pub.Type)

	pub = NewPublisherClient(Mock)
	assert.NotNil(t, pub.Client)
	assert.Equal(t, Mock, pub.Type)

	var activemq BrokerType = "activemq"
	pub = NewPublisherClient(activemq)
	assert.Nil(t, pub.Client)
	assert.Zero(t, pub.Type)
}

func TestNewMessageBrokerClient(t *testing.T) {
	client := NewMessageBrokerClient(RabbitMQ, rabbitmq.NewRMQClient())
	assert.Equal(t, RabbitMQ, client.Type)
	assert.NotNil(t, client.Client)
}

package brokers

import (
	"context"
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

const Mock BrokerType = "mock"

type ClientMock struct {
	mock.Mock
}

func NewMockClient() *ClientMock {
	return new(ClientMock)
}

func OpenMockConnection(client MessageBroker) (*MessageBrokerClient, error) {
	broker = &MessageBrokerClient{Mock, client, 1}
	return broker.connect()
}

func (mock *ClientMock) CloseConnection() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *ClientMock) Connect() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *ClientMock) ConnectionListener() {
	mock.Called()
}

func (mock *ClientMock) Publish(ctx context.Context, queue string, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup) {
	mock.Called(ctx, queue, event, callback, wg)
}

func (mock *ClientMock) PublishMagentoProductEvent(ctx context.Context, queue string, event *base.MagentoProductEvent, callback base.MessageBrokerMagentoProductCallbackFunc, wg *sync.WaitGroup) {
	mock.Called(ctx, queue, event, callback, wg)
}

func (mock *ClientMock) Consume(queue string, callback base.MessageBrokerCallbackFunc) {
	mock.Called(queue, callback)
}

func (mock *ClientMock) ConsumeMagentoProductEvents(queue string, callback base.MessageBrokerMagentoProductCallbackFunc) {
	mock.Called(queue, callback)
}

func (mock *ClientMock) IsConnectionOpen() bool {
	args := mock.Called()
	return args.Bool(0)
}

func (mock *ClientMock) SendToDeadletterQueue(msg amqp.Delivery, err error) error {
	args := mock.Called(msg, err)
	return args.Error(0)
}

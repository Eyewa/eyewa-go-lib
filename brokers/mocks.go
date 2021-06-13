package brokers

import (
	"sync"

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
	broker = &MessageBrokerClient{Mock, client}
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

func (mock *ClientMock) Publish(queue string) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *ClientMock) Consume(wg *sync.WaitGroup, queue string, errChan chan<- error) {
}

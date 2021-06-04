package brokers

import "github.com/stretchr/testify/mock"

const Mock BrokerType = "mock"

type ClientMock struct {
	mock.Mock
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

func (mock *ClientMock) Consume(queue string) error {
	args := mock.Called()
	return args.Error(0)
}

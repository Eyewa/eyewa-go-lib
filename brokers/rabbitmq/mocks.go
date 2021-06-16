package rabbitmq

import (
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/stretchr/testify/mock"
)

type RMQClientMock struct {
	mock.Mock
}

func NewRMQClientMock() *RMQClientMock {
	return new(RMQClientMock)
}

func (mock *RMQClientMock) CloseConnection() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *RMQClientMock) Connect() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *RMQClientMock) Publish(wg *sync.WaitGroup, queue string, eyewaEvent base.EyewaEvent, errChan chan<- error) chan<- error {
	args := mock.Called()
	return args.Get(0).(chan error)
}

func (mock *RMQClientMock) Consume(wg *sync.WaitGroup, queue string, errChan chan<- error) {
}

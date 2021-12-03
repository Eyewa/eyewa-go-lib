package rabbitmq

import (
	"context"
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

func (mock *RMQClientMock) Publish(ctx context.Context, queue string, priority int, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup) {
	_ = mock.Called()
}

func (mock *RMQClientMock) PublishMagentoProductEvent(ctx context.Context, queue string, priority int, event *base.MagentoProductEvent, callback base.MessageBrokerMagentoProductCallbackFunc, wg *sync.WaitGroup) {
	_ = mock.Called()
}

func (mock *RMQClientMock) Consume(queue string, callback base.MessageBrokerCallbackFunc) {
	_ = mock.Called()
}

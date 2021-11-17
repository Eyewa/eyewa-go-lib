package amqp

import (
	"sync"

	"github.com/streadway/amqp"
)

type HeaderCarrier struct {
	mtx sync.Mutex

	gets []string
	sets [][2]string
	data map[string]string
}

func NewHeaderCarrier(data amqp.Table) *HeaderCarrier {
	copied := make(map[string]string, len(data))
	for k, v := range data {
		copied[k] = v.(string)
	}
	return &HeaderCarrier{data: copied}
}

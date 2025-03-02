package domain

import "errors"

var (
	ErrBadQueueSize      = errors.New("max queue size should be greater then zero")
	ErrBadConsumersCount = errors.New("max consumers should be greater then zero")
)

// MQueue очередь сообщений.
type MQueue struct {
	// Name уникальное имя очереди сообщений.
	Name string
}

func NewMQueue(name string, maxSize, maxConsumers int64) (*MQueue, error) {
	if maxSize <= 0 {
		return nil, ErrBadQueueSize
	}

	if maxConsumers <= 0 {
		return nil, ErrBadConsumersCount
	}

	return &MQueue{
		Name: name,
	}, nil
}

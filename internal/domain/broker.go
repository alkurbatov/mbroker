package domain

import (
	"errors"
	"fmt"
)

var (
	ErrQueueExists = errors.New("message queue already exists")
	ErrNoQueue     = errors.New("queue doesn't exist")
)

type Broker struct {
	queues map[string]*MQueue
}

func NewBroker() *Broker {
	return &Broker{
		queues: make(map[string]*MQueue, 0),
	}
}

// RegisterQueue создает новую [Queue] с указанными параметрами.
func (b *Broker) RegisterQueue(name string, maxSize, maxConsumers int) error {
	q, err := NewMQueue(name, maxSize, maxConsumers)
	if err != nil {
		return fmt.Errorf("create message queue '%s': %w", name, err)
	}

	if _, ok := b.queues[name]; ok {
		return ErrQueueExists
	}
	b.queues[name] = q

	return nil
}

// Post размещает сообщение в очереди.
// В случае успеха возвращает количество сообщений, которые можно разместить в очереди.
func (b *Broker) Post(dst string, msg []byte) (int, error) {
	queue := b.queues[dst]
	if queue == nil {
		return 0, ErrNoQueue
	}

	if err := queue.Post(msg); err != nil {
		return 0, fmt.Errorf("post message to queue: %w", err)
	}

	return queue.SpaceLeft(), nil
}

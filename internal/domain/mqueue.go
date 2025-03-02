package domain

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrBadQueueSize      = errors.New("max queue size should be greater then zero")
	ErrBadConsumersCount = errors.New("max consumers should be greater then zero")
)

// MQueue очередь сообщений.
type MQueue struct {
	// Name уникальное имя очереди сообщений.
	Name string

	mu sync.RWMutex

	messages *RingBuffer
}

// NewMQueue создает новую очередь сообщений.
func NewMQueue(name string, maxSize, maxConsumers int) (*MQueue, error) {
	if maxSize <= 0 {
		return nil, ErrBadQueueSize
	}

	if maxConsumers <= 0 {
		return nil, ErrBadConsumersCount
	}

	return &MQueue{
		Name:     name,
		messages: NewRingBuffer(maxSize),
	}, nil
}

// Post размещает сообщение в очереди.
func (q *MQueue) Post(msg Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err := q.messages.PushBack(msg); err != nil {
		return fmt.Errorf("store message: %w", err)
	}

	return nil
}

// SpaceLeft возвращает количество сообщений, которые можно положить в очередь до переполнения.
func (q *MQueue) SpaceLeft() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.messages.SpaceLeft()
}

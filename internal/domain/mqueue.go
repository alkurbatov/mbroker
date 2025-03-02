package domain

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrBadQueueSize      = errors.New("max queue size should be greater then zero")
	ErrBadConsumersCount = errors.New("max consumers should be greater then zero")
	ErrTooManyConsumers  = errors.New("maximum consumers count reached")
	ErrDuplicateConsumer = errors.New("consumer already subscribed")
)

// MQueue очередь сообщений.
type MQueue struct {
	// Name уникальное имя очереди сообщений.
	Name string

	// maxConsumers максимальное допустимое количество подписчиков.
	maxConsumers int

	mu sync.RWMutex

	messages  *RingBuffer
	consumers map[string]struct{}
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
		Name:         name,
		maxConsumers: maxConsumers,
		messages:     NewRingBuffer(maxSize),
		consumers:    make(map[string]struct{}, 0),
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

// AddConsumer добавляет потребителя сообщений.
func (q *MQueue) AddConsumer(clientURL string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.consumers) == q.maxConsumers {
		return ErrTooManyConsumers
	}

	if _, ok := q.consumers[clientURL]; ok {
		return ErrDuplicateConsumer
	}

	q.consumers[clientURL] = struct{}{}

	return nil
}

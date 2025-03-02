package domain

import (
	"errors"
	"sync"
)

var (
	ErrBadConsumersCount = errors.New("max consumers should be greater then zero")
	ErrBadQueueSize      = errors.New("max queue size should be greater then zero")
	ErrDuplicateConsumer = errors.New("consumer already subscribed")
	ErrQueueOverflow     = errors.New("max queue size reached")
	ErrTooManyConsumers  = errors.New("maximum consumers count reached")
)

// MQueue очередь сообщений.
type MQueue struct {
	// Name уникальное имя очереди сообщений.
	Name string

	// maxSize максимальное количество сообщений, ожидающих отправки.
	maxSize int

	// maxConsumers максимальное допустимое количество подписчиков.
	maxConsumers int

	mu sync.RWMutex

	messages  chan Message
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
		maxSize:      maxSize,
		maxConsumers: maxConsumers,
		messages:     make(chan Message, maxSize),
		consumers:    make(map[string]struct{}, 0),
	}, nil
}

// Post размещает сообщение в очереди.
func (q *MQueue) Post(msg Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) == q.maxSize {
		return ErrQueueOverflow
	}

	q.messages <- msg

	return nil
}

// SpaceLeft возвращает количество сообщений, которые можно положить в очередь до переполнения.
func (q *MQueue) SpaceLeft() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.maxSize - len(q.messages)
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

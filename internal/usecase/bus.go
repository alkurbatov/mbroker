package usecase

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/alkurbatov/mbroker/internal/domain"
)

var (
	ErrQueueExists = errors.New("message queue already exists")
	ErrNoQueue     = errors.New("queue doesn't exist")
)

type Bus struct {
	logger *slog.Logger

	queues map[string]*domain.MQueue
}

func NewBus(logger *slog.Logger) Bus {
	return Bus{
		logger: logger,
		queues: make(map[string]*domain.MQueue, 0),
	}
}

// RegisterQueue создает новую [Queue] с указанными параметрами.
func (uc Bus) RegisterQueue(name string, maxSize, maxConsumers int) error {
	q, err := domain.NewMQueue(name, maxSize, maxConsumers)
	if err != nil {
		return fmt.Errorf("create message queue '%s': %w", name, err)
	}

	if _, ok := uc.queues[name]; ok {
		return ErrQueueExists
	}
	uc.queues[name] = q

	return nil
}

// PostMessage размещает сообщение в указанной очереди.
func (uc Bus) PostMessage(queueName string, msg domain.Message) error {
	queue := uc.queues[queueName]
	if queue == nil {
		return ErrNoQueue
	}

	if err := queue.Post(msg); err != nil {
		return fmt.Errorf("post message to queue: %w", err)
	}

	uc.logger.Info("posted message to queue",
		"name", queueName, "space_left", queue.SpaceLeft(), "msg", string(msg))

	return nil
}

// Subscribe подписывает клиента на события очереди.
func (uc Bus) Subscribe(queueName, clientURI string) error {
	queue := uc.queues[queueName]
	if queue == nil {
		return ErrNoQueue
	}

	if err := queue.AddConsumer(clientURI); err != nil {
		return fmt.Errorf("add consumer: %w", err)
	}

	uc.logger.Info("client subscribed to queue",
		"name", queueName, "client", clientURI)

	return nil
}

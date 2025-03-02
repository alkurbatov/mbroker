package usecase

import (
	"fmt"
	"log/slog"

	"github.com/alkurbatov/mbroker/internal/domain"
)

type Producer struct {
	logger *slog.Logger
	broker *domain.Broker
}

func NewProducer(logger *slog.Logger, broker *domain.Broker) Producer {
	return Producer{
		logger: logger,
		broker: broker,
	}
}

// PostMessage размещает сообщение в указанной очереди.
func (uc Producer) PostMessage(dst string, msg domain.Message) error {
	spaceLeft, err := uc.broker.Post(dst, msg)
	if err != nil {
		return fmt.Errorf("post message to broker: %w", err)
	}

	uc.logger.Info("posted message to queue", "name", dst, "space_left", spaceLeft)

	return nil
}

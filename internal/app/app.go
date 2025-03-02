package app

import (
	"fmt"
	"log/slog"

	"github.com/alkurbatov/mbroker/internal/config"
	"github.com/alkurbatov/mbroker/internal/domain"
)

func Run(l *slog.Logger, cfg *config.Config) error {
	queues := make([]*domain.MQueue, 0, len(cfg.Queues))
	for name, settings := range cfg.Queues {
		q, err := domain.NewMQueue(name, settings.MaxSize, settings.MaxConsumers)
		if err != nil {
			return fmt.Errorf("create message queue '%s': %w", name, err)
		}

		queues = append(queues, q)
	}

	for _, q := range queues {
		l.Info("queue created", "name", q.Name)
	}

	return nil
}

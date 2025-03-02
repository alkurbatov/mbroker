package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/alkurbatov/mbroker/internal/api/http/v1"
	"github.com/alkurbatov/mbroker/internal/config"
	"github.com/alkurbatov/mbroker/internal/domain"
	"github.com/alkurbatov/mbroker/internal/usecase"
)

func Run(l *slog.Logger, cfg *config.Config) error {
	broker := domain.NewBroker()

	for name, settings := range cfg.Queues {
		err := broker.RegisterQueue(name, settings.MaxSize, settings.MaxConsumers)
		if err != nil {
			return fmt.Errorf("apply queue settings: %w", err)
		}

		l.Info("queue created", "name", name, "max_size", settings.MaxSize)
	}

	producer := usecase.NewProducer(l, broker)

	router := gin.Default()
	v1.Inject(router, producer)

	if err := router.Run(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run HTTP API: %w", err)
	}

	return nil
}

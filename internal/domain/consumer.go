package domain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alkurbatov/mbroker/internal/infra/logging"
)

var ErrNotificationFailed = errors.New("notification failed")

type Consumer struct {
	uri string
	src <-chan Message
}

// NewConsumer создает нового подписчика.
func NewConsumer(uri string, src <-chan Message) *Consumer {
	c := &Consumer{
		uri: uri,
		src: src,
	}

	go c.listen()

	return c
}

// Post отправляет сообщение подписчику.
func (c *Consumer) Post(msg Message) error {
	body := bytes.NewBuffer(msg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.uri,
		body,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: code %d", ErrNotificationFailed, resp.StatusCode)
	}

	return nil
}

func (c *Consumer) listen() {
	l := slog.Default()

	for {
		msg := <-c.src

		if err := c.Post(msg); err != nil {
			l.Error("notification failed", logging.Err(err))
			continue
		}

		l.Info("notification sent", "consumer", c.uri, "msg", string(msg))
	}
}

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

var ErrRequestFailed = errors.New("request failure")

func listen(router http.Handler, address string) {
	httpSrv := &http.Server{
		Handler:           router,
		Addr:              address,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func subscribe(address string) error {
	jsonStr := fmt.Sprintf(`{"uri":"http://%s"}`, address)
	body := bytes.NewBufferString(jsonStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://0.0.0.0:8080/v1/queues/mail/subscriptions",
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
		return fmt.Errorf("%w: code %s", ErrRequestFailed, resp.Status)
	}

	return nil
}

func main() {
	l := slog.Default()

	address := fmt.Sprintf(
		"0.0.0.0:%d",
		9000+rand.Intn(1000), //nolint: gosec //ОК для тестовой утилиты
	)

	router := http.NewServeMux()
	router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			l.Error("parse message", "error", err)
			http.Error(w, "parse message", http.StatusBadRequest)
			return
		}

		l.Info("Received message", "msg", string(body))
	})

	go listen(router, address)

	if err := subscribe(address); err != nil {
		l.Error("subscribe to queue", "error", err)
		return
	}

	l.Info("subscribed")

	select {}
}

package main

import (
	"log/slog"
	"os"

	"github.com/alkurbatov/mbroker/internal/app"
	"github.com/alkurbatov/mbroker/internal/config"
	"github.com/alkurbatov/mbroker/internal/infra/logging"
)

func run() int {
	l := slog.Default()

	cfg, err := config.New()
	if err != nil {
		l.Error("new config", logging.Err(err))
		return 1
	}

	if err = app.Run(l, cfg); err != nil {
		l.Error("run service", logging.Err(err))
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}

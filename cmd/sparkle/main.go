package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/aqyuki/sparkle/internal/bot/handler"
	"github.com/aqyuki/sparkle/internal/di"
	"github.com/aqyuki/sparkle/pkg/env"
	"github.com/aqyuki/sparkle/pkg/logging"
	"github.com/samber/do"
)

type exitCode int

const (
	ExitCodeOK exitCode = iota
	ExitCodeError
)

func main() {
	logger := logging.NewLoggerFromEnv()
	ctx := logging.WithLogger(context.Background(), logger)
	defer exit(run(ctx))
}

func run(ctx context.Context) exitCode {
	logger := logging.FromContext(ctx)
	defer logger.Sync()

	logger.Info("loading configuration")
	token, err := env.GetOrErr("DISCORD_TOKEN")
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		return ExitCodeError
	}
	logger.Info("loaded configuration")

	// init ID container
	injector := do.New()
	do.Provide(injector, di.NewLoggerInjector(logger))
	do.Provide(injector, di.NewSessionInjector(token))
	do.Provide(injector, handler.NewReadyHandler)
	do.Provide(injector, handler.NewMessageLinkExpandHandler)
	do.Provide(injector, bot.NewBot)

	// health check
	if err := do.HealthCheck[*bot.Bot](injector); err != nil {
		logger.Error(err.Error(), "error", err)
		return ExitCodeError
	}

	// get b
	logger.Info("initializing bot")
	b, err := do.Invoke[*bot.Bot](injector)
	if err != nil {
		logger.Error(err.Error(), "error", err)
		return ExitCodeError
	}

	logger.Info("starting bot")
	if err := b.Start(); err != nil {
		logger.Error(err.Error(), "error", err)
		return ExitCodeError
	}

	// wait for shutdown signal
	logger.Info("waiting for shutdown signal")
	ctx, done := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer done()
	<-ctx.Done()
	logger.Info("received shutdown signal")

	// shutdown
	logger.Info("shutting down")
	if err := do.Shutdown[*bot.Bot](injector); err != nil {
		logger.Error(err.Error(), "error", err)
		return ExitCodeError
	}

	logger.Info("shutdown complete. goodbye!")
	return ExitCodeOK
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}

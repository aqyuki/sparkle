package main

import (
	"context"
	"os"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/aqyuki/sparkle/pkg/env"
	"github.com/aqyuki/sparkle/pkg/logging"
	"go.uber.org/zap"
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

	// initialize configuration
	logger.Info("configuration is loading")
	token, err := env.GetOrErr("DISCORD_TOKEN")
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		return ExitCodeError
	}
	logger.Info("configuration was loaded successfully")

	logger.Infof("bot initialization is starting")
	b, err := bot.New(token,
		&bot.Deps{
			Logger: logger,
		})
	if err != nil {
		logger.Error("failed to initialize bot", "error", err)
		return ExitCodeError
	}

	if err := b.Run(token); err != nil {
		logger.Error("failed to run bot", "error", err)
		return ExitCodeError
	}

	<-ctx.Done()
	logger.Infof("shutdown signal was received", zap.String("reason", ctx.Err().Error()))
	if err := b.Shutdown(); err != nil {
		logger.Error("failed to shutdown bot", "error", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}

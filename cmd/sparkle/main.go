package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/aqyuki/sparkle/internal/bot/command"
	"github.com/aqyuki/sparkle/internal/bot/handler"
	"github.com/aqyuki/sparkle/internal/information"
	"github.com/aqyuki/sparkle/pkg/cache"
	"github.com/aqyuki/sparkle/pkg/env"
	"github.com/aqyuki/sparkle/pkg/logging"
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

	// initialize handlers

	// Ready handler
	infoProvider := information.NewBotInformationProvider()
	readyHandler := handler.NewReadyHandler(logger, infoProvider)

	// MessageCreate handler
	cache := cache.NewImMemoryCacheStore(5*time.Minute, 10*time.Minute)
	msgLinkExpandHandler := handler.NewMessageLinkExpandHandler(logger, cache)

	// Message Command Runner
	mentionRouter := command.NewMentionCommandRouter()

	// register commands
	versionCmd := command.NewVersionCommand(logger, infoProvider)
	mentionRouter.Register(versionCmd)

	// initialize bot
	logger.Info("initializing bot")
	bot, err := bot.NewBot(token)
	if err != nil {
		logger.Error(err.Error(), "error", err)
		return ExitCodeError
	}

	// register handlers
	bot.AddReadyHandler(readyHandler)
	bot.AddMessageCreateHandler(msgLinkExpandHandler)
	bot.AddMessageCreateHandler(mentionRouter)

	logger.Info("starting bot")
	if err := bot.Start(); err != nil {
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
	if err := bot.Shutdown(); err != nil {
		logger.Error(err.Error(), "error", err)
		return ExitCodeError
	}
	logger.Info("shutdown complete. goodbye!")
	return ExitCodeOK
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}

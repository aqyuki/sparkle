package bot

import (
	"time"

	"github.com/aqyuki/sparkle/internal/bot/command"
	"github.com/aqyuki/sparkle/internal/bot/handler"
	"github.com/aqyuki/sparkle/internal/bot/internal/core"
	"github.com/aqyuki/sparkle/internal/bot/router"
	"github.com/aqyuki/sparkle/internal/information"
	"github.com/aqyuki/sparkle/pkg/cache"
	"go.uber.org/zap"
)

// Bot is a struct to provide bot features.
type Bot struct {
	core   *core.Core
	logger *zap.SugaredLogger
	info   *information.BotInformation
}

func New(token string, deps *Deps) (*Bot, error) {
	c, err := core.New(token)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		core:   c,
		logger: deps.Logger,
		info:   information.NewBotInformation(),
	}

	// initialize handler
	ready := handler.NewReadyHandler(b.logger, b.info)
	expandMessage := handler.NewMessageLinkExpandHandler(b.logger, cache.NewImMemoryCacheStore(5*time.Minute, 10*time.Minute))

	// initialize router
	commandRouter := router.NewCommandRouter("s!", b.logger)
	commandRouter.Register("version", command.NewVersionCommand(b.logger, b.info))

	// register handlers
	b.core.AddReadyHandler(ready.HandleReady)
	b.core.AddMessageCreateHandler(expandMessage.Expand)
	b.core.AddMessageCreateHandler(commandRouter.Handle)
	return b, nil
}

func (b *Bot) Run(token string) error {
	b.logger.Infof("bot is starting...")
	return b.core.Open(token)
}

func (b *Bot) Shutdown() error {
	b.logger.Infof("bot is shutting down...")
	return b.core.Close()
}

type Deps struct {
	Logger *zap.SugaredLogger
}

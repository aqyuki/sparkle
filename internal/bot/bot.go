package bot

import (
	"github.com/aqyuki/sparkle/internal/bot/handler"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/do"
	"github.com/samber/oops"
)

var _ do.Shutdownable = (*Bot)(nil)
var _ do.Provider[*Bot] = NewBot
var _ do.Healthcheckable = (*Bot)(nil)

type Bot struct {
	session   *discordgo.Session
	ready     handler.ReadyHandler
	msgExpand handler.MessageLinkExpandHandler
	remover   []func()
}

// NewBot creates a new Bot instance.
// This function is used by the dependency injector to create a new Bot instance.
func NewBot(i *do.Injector) (*Bot, error) {
	// Resolve dependencies
	session, err := do.Invoke[*discordgo.Session](i)
	if err != nil {
		return nil, oops.
			In("Bot").
			Errorf("dependency resolution failed: %w", err)
	}

	ready, err := do.Invoke[handler.ReadyHandler](i)
	if err != nil {
		return nil, oops.
			In("Bot").
			Errorf("dependency resolution failed: %w", err)
	}

	msgExpand, err := do.Invoke[handler.MessageLinkExpandHandler](i)
	if err != nil {
		return nil, oops.
			In("Bot").
			Errorf("dependency resolution failed: %w", err)
	}

	return &Bot{
		session:   session,
		ready:     ready,
		msgExpand: msgExpand,
		remover:   make([]func(), 0),
	}, nil
}

func (b *Bot) Start() error {
	b.remover = append(b.remover,
		b.session.AddHandler(b.ready.Ready),
	)
	if err := b.session.Open(); err != nil {
		return oops.
			In("Bot").
			Errorf("failed to open session: %w", err)
	}
	return nil
}

func (b *Bot) Shutdown() error {
	for _, f := range b.remover {
		f()
	}
	if err := b.session.Close(); err != nil {
		return oops.
			In("Bot").
			Errorf("failed to close session: %w", err)
	}
	return nil
}

func (b *Bot) HealthCheck() error {
	if b.session == nil {
		return oops.
			In("Bot").
			Errorf("session is not ready")
	}
	return nil
}

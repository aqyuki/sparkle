package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/samber/oops"
)

type Bot struct {
	session *discordgo.Session
	remover []func()
}

// NewBot creates a new Bot instance.
// This function is used by the dependency injector to create a new Bot instance.
func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, oops.
			In("Bot").
			Errorf("failed to create session: %w", err)
	}
	return &Bot{
		session: session,
		remover: make([]func(), 0),
	}, nil
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return oops.
			In("Bot").
			Errorf("failed to open session: %w", err)
	}
	return nil
}

func (b *Bot) AddReadyHandler(handler ReadyHandler) {
	f := b.session.AddHandler(handler.Handle)
	b.remover = append(b.remover, f)
}

func (b *Bot) AddMessageCreateHandler(handler MessageCreateHandler) {
	f := b.session.AddHandler(handler.Handle)
	b.remover = append(b.remover, f)
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

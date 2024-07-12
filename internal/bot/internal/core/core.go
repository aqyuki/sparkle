package core

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/oops"
)

type Core struct {
	session *discordgo.Session
	remover []func()
}

func New(token string) (*Core, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, oops.
			In("Core").
			Time(time.Now()).
			Wrapf(err, "failed to create session")
	}

	return &Core{
		remover: make([]func(), 0),
		session: session,
	}, nil
}

// Open opens a new session with the given token.
func (c *Core) Open(token string) error {

	if err := c.session.Open(); err != nil {
		return oops.
			In("Core").
			Time(time.Now()).
			Wrapf(err, "failed to open session")
	}
	return nil
}

// Close closes the session.
func (c *Core) Close() error {
	// if session is nil, return an error
	if c.session == nil {
		return oops.
			In("Core").
			Time(time.Now()).
			Errorf("session is not open")
	}

	// remove handlers
	for i := range c.remover {
		c.remover[i]()
	}

	if err := c.session.Close(); err != nil {
		return oops.
			In("Core").
			Time(time.Now()).
			Wrapf(err, "failed to close session")
	}
	return nil
}

func (c *Core) addHandler(handler any) {
	f := c.session.AddHandler(handler)
	c.remover = append(c.remover, f)
}

func (c *Core) AddReadyHandler(handler func(*discordgo.Session, *discordgo.Ready)) {
	c.addHandler(handler)
}

func (c *Core) AddMessageCreateHandler(handler func(*discordgo.Session, *discordgo.MessageCreate)) {
	c.addHandler(handler)
}

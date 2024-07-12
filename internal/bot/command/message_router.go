package command

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/oops"
)

type MessageRouter struct {
	prefix   string
	commands map[string]Command
}

func NewMessageRouter(prefix string) *MessageRouter {
	return &MessageRouter{
		prefix:   prefix,
		commands: make(map[string]Command),
	}
}

func (r *MessageRouter) Register(cmd Command) error {
	if _, ok := r.commands[strings.ToLower(cmd.Name())]; ok {
		return oops.Errorf("command %s is already registered", cmd.Name())
	}
	r.commands[strings.ToLower(cmd.Name())] = cmd
	return nil
}

func (r *MessageRouter) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !strings.HasPrefix(message.Content, r.prefix) {
		return
	}
	for k, v := range r.commands {
		if !strings.Contains(message.Content, k) {
			continue
		}
		v.Handle(session, message)
	}
}

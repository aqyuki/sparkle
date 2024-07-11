package command

import (
	"strings"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/oops"
)

type CommandRouter interface{ bot.MessageCreateHandler }

type Command interface {
	Name() string
	Handle(*discordgo.Session, *discordgo.MessageCreate)
}

type MentionCommandRouter struct {
	commands map[string]Command
}

func NewMentionCommandRouter() *MentionCommandRouter {
	return &MentionCommandRouter{
		commands: make(map[string]Command),
	}
}

func (r *MentionCommandRouter) Register(cmd Command) error {
	if _, ok := r.commands[strings.ToLower(cmd.Name())]; ok {
		return oops.Errorf("command %s is already registered", cmd.Name())
	}
	r.commands[strings.ToLower(cmd.Name())] = cmd
	return nil
}

func (r *MentionCommandRouter) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	for _, mention := range message.Mentions {
		if mention.ID != session.State.User.ID {
			continue
		}
	}
	for k, v := range r.commands {
		if !strings.Contains(message.Content, k) {
			continue
		}
		v.Handle(session, message)
	}
}

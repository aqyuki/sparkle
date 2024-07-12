package router

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

type Command interface {
	Run(session *discordgo.Session, message *discordgo.MessageCreate, args []string)
}

type CommandRouter struct {
	prefix   string
	logger   *zap.SugaredLogger
	commands map[string]Command
}

func NewCommandRouter(prefix string, logger *zap.SugaredLogger) *CommandRouter {
	return &CommandRouter{
		prefix:   prefix,
		logger:   logger,
		commands: make(map[string]Command),
	}
}

func (r *CommandRouter) Register(name string, command Command) error {
	if command == nil {
		return oops.
			In("Command Router").
			Time(time.Now()).
			Errorf("command is nil")
	}
	if _, ok := r.commands[name]; ok {
		return oops.
			In("Command Router").
			Time(time.Now()).
			Errorf("command %s is already registered", name)
	}
	r.commands[name] = command
	return nil
}

func (r *CommandRouter) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		r.logger.Info("bot message was skipped")
		return
	}

	content := strings.Split(message.Content, " ")
	if len(content) == 0 {
		r.logger.Infof("empty message was skipped")
		return
	}

	if !strings.HasPrefix(content[0], r.prefix) {
		r.logger.Infof("message was not a command")
		return
	}

	if cmd, ok := r.commands[strings.TrimPrefix(content[0], r.prefix)]; ok {
		cmd.Run(session, message, content[1:])
		return
	}
}

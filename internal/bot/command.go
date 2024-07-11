package bot

import (
	"github.com/bwmarrin/discordgo"
)

var _ CommandRunner = (*MessageCommandRunner)(nil)

type Command interface {
	Handle(session *discordgo.Session, message *discordgo.MessageCreate)
}

type CommandRunner interface{ MessageCreateHandler }

type MessageCommandRunner struct {
	commands []Command
}

func NewMessageCommandRunner() *MessageCommandRunner {
	return &MessageCommandRunner{}
}

func (m *MessageCommandRunner) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	var contained bool
	for _, mention := range message.Mentions {
		if mention.ID == session.State.User.ID {
			contained = true
		}
	}
	if !contained {
		return
	}
	for _, command := range m.commands {
		command.Handle(session, message)
	}
}

func (m *MessageCommandRunner) RegisterCommand(command Command) {
	m.commands = append(m.commands, command)
}

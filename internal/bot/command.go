package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var _ CommandRunner = (*MessageCommandRunner)(nil)

type Command interface {
	Handle(session *discordgo.Session, message *discordgo.MessageCreate)
}

type CommandRunner interface{ MessageCreateHandler }

type MessageCommandRunner struct {
	option   *CommandRunnerOption
	commands []Command
}

func NewMessageCommandRunner(option *CommandRunnerOption) *MessageCommandRunner {
	if option == nil {
		option = NewDefaultCommandRunnerOption()
	}
	return &MessageCommandRunner{option: option}
}

func (m *MessageCommandRunner) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}
	if !strings.HasPrefix(message.Content, m.option.CommandPrefix()) {
		return
	}
	for _, command := range m.commands {
		command.Handle(session, message)
	}
}

func (m *MessageCommandRunner) RegisterCommand(command Command) {
	m.commands = append(m.commands, command)
}

type CommandRunnerOption struct {
	Prefix    string
	Separator string
}

func NewCommandRunnerOption(ops ...OptionProperty) *CommandRunnerOption {
	option := NewDefaultCommandRunnerOption()
	for _, op := range ops {
		op(option)
	}
	return option
}

func NewDefaultCommandRunnerOption() *CommandRunnerOption {
	return &CommandRunnerOption{
		Prefix:    "spk",
		Separator: "!",
	}
}

func (o *CommandRunnerOption) CommandPrefix() string {
	return fmt.Sprintf("%s%s", o.Prefix, o.Separator)
}

type OptionProperty func(*CommandRunnerOption)

func WithPrefix(prefix string) OptionProperty {
	return func(o *CommandRunnerOption) {
		o.Prefix = prefix
	}
}

func WithSeparator(separator string) OptionProperty {
	return func(o *CommandRunnerOption) {
		o.Separator = separator
	}
}

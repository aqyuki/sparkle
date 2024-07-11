package command

import (
	"fmt"
	"strings"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/aqyuki/sparkle/internal/information"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ bot.Command = (*VersionCommand)(nil)

type VersionCommand struct {
	logger      *zap.SugaredLogger
	information information.InformationProvider
}

func NewVersionCommand(logger *zap.SugaredLogger, info information.InformationProvider) *VersionCommand {
	return &VersionCommand{
		logger:      logger,
		information: info,
	}
}

func (c *VersionCommand) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	if !strings.Contains(message.Content, "version") {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Botのバージョン",
		Description: fmt.Sprintf("現在のBotのバージョンは `%s` です", c.information.Version()),
		Color:       0x7fffff,
	}

	msg := &discordgo.MessageSend{
		Embeds:          []*discordgo.MessageEmbed{embed},
		AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
		Reference:       message.Reference(),
	}
	if _, err := session.ChannelMessageSendComplex(message.ChannelID, msg); err != nil {
		c.logger.Errorf("failed to send message: %v", err)
	}
}

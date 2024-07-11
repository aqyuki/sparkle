package command

import (
	"fmt"

	"github.com/aqyuki/sparkle/internal/information"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ Command = (*VersionCommand)(nil)

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

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	if _, err := session.ChannelMessageSendComplex(message.ChannelID,
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "バージョン情報",
					Description: fmt.Sprintf("現在のBotのバージョンは `%s` です。", c.information.Version()),
					Color:       0x7fffff,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
			Reference:       message.Reference(),
		}); err != nil {
		c.logger.Errorf("failed to send message: %v", err)
	}
}

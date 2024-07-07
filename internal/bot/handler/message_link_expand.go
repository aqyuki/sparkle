package handler

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/aqyuki/sparkle/pkg/logging"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/do"
	"go.uber.org/zap"
)

var _ MessageLinkExpandHandler = (*messageLinkExpandHandler)(nil)
var _ do.Provider[MessageLinkExpandHandler] = NewMessageLinkExpandHandler

type MessageLinkExpandHandler interface {
	Expand(session *discordgo.Session, message *discordgo.MessageCreate)
}

type messageLinkExpandHandler struct {
	logger *zap.SugaredLogger
	rgx    *regexp.Regexp
}

func NewMessageLinkExpandHandler(i *do.Injector) (MessageLinkExpandHandler, error) {
	logger, err := do.Invoke[*zap.SugaredLogger](i)
	if err != nil {
		logger = logging.DefaultLogger()
		logger.Warn("dependency resolution failed for *zap.SugaredLogger and recovered with the default logger")
	}

	return &messageLinkExpandHandler{
		logger: logger,
		rgx:    regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord(app)?\.com/channels/(\d+)/(\d+)/(\d+)`),
	}, nil
}

func (h *messageLinkExpandHandler) Expand(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		h.logger.Info("skip bot message")
		return
	}

	links := h.extractLink(message.Content)
	if len(links) == 0 {
		h.logger.Info("no message link found")
		return
	}

	// 荒らし対策で，リンクが複数ある場合は最初のリンクのみ展開する
	link := links[0]
	info, err := h.extractMessageInfo(link)
	if err != nil {
		h.logger.Errorf("failed to extract message info: %v", err)
		return
	}

	// 違うギルドのメッセージは展開しない
	if info.guild != message.GuildID {
		h.logger.Info("skip message from different guild")
		return
	}

	// 対象のメッセージが投稿されたチャンネルがNSFWの場合は展開しない
	channel, err := session.Channel(info.channel)
	if err != nil {
		h.logger.Errorf("failed to get channel: %v", err)
		return
	}
	if channel.NSFW {
		h.logger.Info("skip NSFW channel")
		return
	}

	// 対象メッセージを取得
	msg, err := session.ChannelMessage(info.channel, info.message)
	if err != nil {
		h.logger.Errorf("failed to get message: %v", err)
		return
	}

	var image *discordgo.MessageEmbedImage
	if len(msg.Attachments) > 0 {
		image = &discordgo.MessageEmbedImage{
			URL: msg.Attachments[0].URL,
		}
	}

	embed := &discordgo.MessageEmbed{
		Image: image,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    msg.Author.Username,
			IconURL: msg.Author.AvatarURL("64"),
		},
		Color:       0x7fffff,
		Description: msg.Content,
		Timestamp:   msg.Timestamp.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: channel.Name,
		},
	}

	replyMsg := discordgo.MessageSend{
		Embeds:    []*discordgo.MessageEmbed{embed},
		Reference: message.Reference(),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			RepliedUser: true,
		},
	}

	if _, err := session.ChannelMessageSendComplex(message.ChannelID, &replyMsg); err != nil {
		h.logger.Errorf("failed to send message: %v", err)
		return
	}
	h.logger.Info("message link expanded")
}

func (h *messageLinkExpandHandler) extractLink(content string) []string {
	return h.rgx.FindAllString(content, -1)
}

// extractMessageInfo extracts the channel ID and message ID from the message link.
func (h *messageLinkExpandHandler) extractMessageInfo(link string) (info message, err error) {
	segments := strings.Split(link, "/")
	if len(segments) < 4 {
		return message{}, errors.New("invalid message link")
	}
	return message{
		guild:   segments[len(segments)-3],
		channel: segments[len(segments)-2],
		message: segments[len(segments)-1],
	}, nil
}

type message struct {
	guild   string
	channel string
	message string
}

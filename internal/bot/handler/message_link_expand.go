package handler

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/aqyuki/sparkle/pkg/cache"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ bot.MessageCreateHandler = (*MessageLinkExpandHandler)(nil)

type MessageLinkExpandHandler struct {
	logger *zap.SugaredLogger
	rgx    *regexp.Regexp
	cache  cache.CacheStore
}

func NewMessageLinkExpandHandler(logger *zap.SugaredLogger, cache cache.CacheStore) *MessageLinkExpandHandler {
	return &MessageLinkExpandHandler{
		logger: logger,
		rgx:    regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord(app)?\.com/channels/(\d+)/(\d+)/(\d+)`),
		cache:  cache,
	}
}

func (h *MessageLinkExpandHandler) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
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
	// Cacheを検索する
	var channel *discordgo.Channel
	if v, ok := h.cache.Get(info.channel); ok {
		channel = v.(*discordgo.Channel)
	} else {
		// Cacheに存在しない場合，APIから取得する
		ch, err := session.Channel(info.channel)
		if err != nil {
			h.logger.Errorf("failed to get channel: %v", err)
			return
		}
		channel = ch
		h.cache.Set(info.channel, channel)
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

	// メッセージにリアクションが有る場合は，Embedのフィールドに追加する．
	var field *discordgo.MessageEmbedField
	if len(msg.Reactions) > 0 {
		field = &discordgo.MessageEmbedField{
			Name:   "Reactions",
			Value:  "",
			Inline: true,
		}
		for _, reaction := range msg.Reactions {
			var emoji string
			if reaction.Emoji.ID != "" {
				raw, err := session.GuildEmoji(channel.GuildID, reaction.Emoji.ID)
				if err != nil {
					h.logger.Errorf("failed to get emoji: %v", err)
					continue
				}
				emoji = raw.MessageFormat()
			} else {
				emoji = reaction.Emoji.Name
			}
			field.Value = fmt.Sprintf("%s%s ", field.Value, emoji)
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
		Fields:      fieldOrNil(field),
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

func (h *MessageLinkExpandHandler) extractLink(content string) []string {
	return h.rgx.FindAllString(content, -1)
}

// extractMessageInfo extracts the channel ID and message ID from the message link.
func (h *MessageLinkExpandHandler) extractMessageInfo(link string) (info message, err error) {
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

func fieldOrNil(field *discordgo.MessageEmbedField) []*discordgo.MessageEmbedField {
	if field == nil {
		return nil
	}
	return []*discordgo.MessageEmbedField{field}
}

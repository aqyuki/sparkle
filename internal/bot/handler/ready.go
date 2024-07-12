package handler

import (
	"github.com/aqyuki/sparkle/internal/information"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type ReadyHandler struct {
	logger *zap.SugaredLogger
	info   *information.BotInformation
}

func NewReadyHandler(logger *zap.SugaredLogger, info *information.BotInformation) *ReadyHandler {
	return &ReadyHandler{
		logger: logger,
		info:   info,
	}
}

func (h *ReadyHandler) HandleReady(session *discordgo.Session, event *discordgo.Ready) {
	h.logger.Infof("bot is ready as %s#%s", event.User.Username, event.User.Discriminator)
	h.logger.Infof("version: %s", h.info.Version)
}

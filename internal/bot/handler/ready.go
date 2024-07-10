package handler

import (
	"fmt"

	"github.com/aqyuki/sparkle/internal/bot"
	"github.com/aqyuki/sparkle/internal/information"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ bot.ReadyHandler = (*ReadyHandler)(nil)

type ReadyHandler struct {
	logger          *zap.SugaredLogger
	versionResolver information.InformationProvider
}

func NewReadyHandler(logger *zap.SugaredLogger, resolver information.InformationProvider) *ReadyHandler {
	return &ReadyHandler{
		logger:          logger,
		versionResolver: resolver,
	}
}

func (h *ReadyHandler) Handle(session *discordgo.Session, event *discordgo.Ready) {
	if err := session.UpdateCustomStatus("waking up"); err != nil {
		h.logger.Errorf("failed to update custom status: %v", err)
	}
	h.logger.Infof("bot is ready as %s#%s", event.User.Username, event.User.Discriminator)
	h.logger.Infof("version: %s", h.versionResolver.Version())
	if err := session.UpdateCustomStatus(fmt.Sprintf("version : %s", h.versionResolver.Version())); err != nil {
		h.logger.Errorf("failed to update custom status: %v", err)
	}
}

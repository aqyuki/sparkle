package handler

import (
	"github.com/aqyuki/sparkle/pkg/logging"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/do"
	"go.uber.org/zap"
)

var _ ReadyHandler = (*readyHandler)(nil)
var _ do.Provider[ReadyHandler] = NewReadyHandler

type ReadyHandler interface {
	Ready(session *discordgo.Session, event *discordgo.Ready)
}

type readyHandler struct {
	logger *zap.SugaredLogger
}

func NewReadyHandler(i *do.Injector) (ReadyHandler, error) {
	logger, err := do.Invoke[*zap.SugaredLogger](i)
	if err != nil {
		logger = logging.DefaultLogger()
		logger.Warn("dependency resolution failed for *zap.SugaredLogger and recovered with the default logger")
	}
	return &readyHandler{
		logger: logger,
	}, nil
}

func (h *readyHandler) Ready(session *discordgo.Session, event *discordgo.Ready) {
	if err := session.UpdateCustomStatus("waking up"); err != nil {
		h.logger.Errorf("failed to update custom status: %v", err)
	}
	h.logger.Infof("bot is ready as %s#%s", event.User.Username, event.User.Discriminator)
	if err := session.UpdateCustomStatus("waiting any action"); err != nil {
		h.logger.Errorf("failed to update custom status: %v", err)
	}
}

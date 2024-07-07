package di

import (
	"github.com/bwmarrin/discordgo"
	"github.com/samber/do"
	"go.uber.org/zap"
)

func NewSessionInjector(token string) do.Provider[*discordgo.Session] {
	return func(i *do.Injector) (*discordgo.Session, error) {
		return discordgo.New("Bot " + token)
	}
}

func NewLoggerInjector(logger *zap.SugaredLogger) do.Provider[*zap.SugaredLogger] {
	return func(i *do.Injector) (*zap.SugaredLogger, error) {
		return logger, nil
	}
}

package bot

import (
	"testing"

	"github.com/aqyuki/sparkle/internal/bot/handler"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	t.Parallel()

	type deps struct {
		session   do.Provider[*discordgo.Session]
		ready     do.Provider[handler.ReadyHandler]
		msgExpand do.Provider[handler.MessageLinkExpandHandler]
	}

	tests := []struct {
		name    string
		deps    *deps
		wantErr bool
	}{
		{
			name: "success to create a new bot",
			deps: &deps{
				session: func(i *do.Injector) (*discordgo.Session, error) {
					return &discordgo.Session{}, nil
				},
				ready: func(i *do.Injector) (handler.ReadyHandler, error) {
					return &struct{ handler.ReadyHandler }{}, nil
				},
				msgExpand: func(i *do.Injector) (handler.MessageLinkExpandHandler, error) {
					return &struct {
						handler.MessageLinkExpandHandler
					}{}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "failed to create a new bot because failed to resolve session",
			deps: &deps{
				session: func(i *do.Injector) (*discordgo.Session, error) {
					return nil, assert.AnError
				},
				ready: func(i *do.Injector) (handler.ReadyHandler, error) {
					return &struct{ handler.ReadyHandler }{}, nil
				},
				msgExpand: func(i *do.Injector) (handler.MessageLinkExpandHandler, error) {
					return &struct {
						handler.MessageLinkExpandHandler
					}{}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "failed to create a new bot because failed to resolve ready handler",
			deps: &deps{
				session: func(i *do.Injector) (*discordgo.Session, error) {
					return &discordgo.Session{}, nil
				},
				ready: func(i *do.Injector) (handler.ReadyHandler, error) {
					return nil, assert.AnError
				},
				msgExpand: func(i *do.Injector) (handler.MessageLinkExpandHandler, error) {
					return &struct {
						handler.MessageLinkExpandHandler
					}{}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "failed to create a new bot because failed to resolve message_link_expand",
			deps: &deps{
				session: func(i *do.Injector) (*discordgo.Session, error) {
					return &discordgo.Session{}, nil
				},
				ready: func(i *do.Injector) (handler.ReadyHandler, error) {
					return &struct{ handler.ReadyHandler }{}, nil
				},
				msgExpand: func(i *do.Injector) (handler.MessageLinkExpandHandler, error) {
					return nil, assert.AnError
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup injector
			i := do.New()
			do.Provide(i, tt.deps.session)
			do.Provide(i, tt.deps.ready)
			do.Provide(i, tt.deps.msgExpand)

			got, err := NewBot(i)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

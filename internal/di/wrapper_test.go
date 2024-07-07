package di

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewSessionInjector(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "success to create a new session injector",
			token: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSessionInjector(tt.token)
			if !assert.NotNil(t, actual) {
				return
			}

			i := do.New()
			do.Provide(i, actual)
			session, err := do.Invoke[*discordgo.Session](i)
			assert.NoError(t, err)
			assert.NotNil(t, session)
		})
	}
}

func TestNewLoggerInjector(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		logger *zap.SugaredLogger
	}{
		{
			name:   "success to create a new logger injector",
			logger: zap.L().Sugar(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLoggerInjector(tt.logger)
			assert.NotNil(t, actual)

			i := do.New()
			do.Provide(i, actual)
			session, err := do.Invoke[*zap.SugaredLogger](i)
			assert.NoError(t, err)
			assert.NotNil(t, session)
		})
	}
}

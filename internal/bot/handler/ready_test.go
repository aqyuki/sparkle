package handler

import (
	"testing"

	"github.com/aqyuki/sparkle/internal/info"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewReadyHandler(t *testing.T) {
	type deps struct {
		logger  do.Provider[*zap.SugaredLogger]
		version do.Provider[info.VersionResolver]
	}
	tests := []struct {
		name string
		deps deps
	}{
		{
			name: "success to create a new ready handler",
			deps: deps{
				logger: func(i *do.Injector) (*zap.SugaredLogger, error) {
					return zap.NewExample().Sugar(), nil
				},
				version: func(i *do.Injector) (info.VersionResolver, error) {
					return info.New(i)
				},
			},
		},
		{
			name: "success to create a new ready handler with default logger",
			deps: deps{
				logger: func(i *do.Injector) (*zap.SugaredLogger, error) {
					return nil, assert.AnError
				},
				version: func(i *do.Injector) (info.VersionResolver, error) {
					return info.New(i)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init DI container
			i := do.New()
			do.Provide(i, tt.deps.logger)
			do.Provide(i, tt.deps.version)

			actual, err := NewReadyHandler(i)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
		})
	}
}

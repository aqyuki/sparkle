package handler

import (
	"testing"

	"github.com/aqyuki/sparkle/pkg/cache"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewMessageLinkExpandHandler(t *testing.T) {
	type deps struct {
		logger do.Provider[*zap.SugaredLogger]
		cache  do.Provider[cache.CacheStore]
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
				cache: func(i *do.Injector) (cache.CacheStore, error) {
					return &struct{ cache.CacheStore }{}, nil
				},
			},
		},
		{
			name: "success to create a new ready handler with default logger",
			deps: deps{
				logger: func(i *do.Injector) (*zap.SugaredLogger, error) {
					return nil, assert.AnError
				},
				cache: func(i *do.Injector) (cache.CacheStore, error) {
					return &struct{ cache.CacheStore }{}, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init DI container
			i := do.New()
			do.Provide(i, tt.deps.logger)
			do.Provide(i, tt.deps.cache)

			actual, err := NewMessageLinkExpandHandler(i)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
		})
	}
}

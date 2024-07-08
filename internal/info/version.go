package info

import "github.com/samber/do"

const defaultVersion = "unknown"

var (
	// Version is the version of the application.
	// This is set at build time using the -ldflags "-X 'internal/info.Version=$VERSION'"
	Version = defaultVersion
)

type VersionResolver interface {
	Version() string
}

type versionResolver struct{}

func New(_ *do.Injector) (VersionResolver, error) {
	return &versionResolver{}, nil
}

func (r *versionResolver) Version() string {
	return Version
}

package information

const defaultVersion = "unknown"

var _ InformationProvider = (*BotInformationProvider)(nil)

// Version is the version of the application.
// This is set at build time using the -ldflags "-X 'internal/info.Version=$VERSION'"
var Version = defaultVersion

type InformationProvider interface {
	Version() string
}

type BotInformationProvider struct{}

func NewBotInformationProvider() *BotInformationProvider {
	return &BotInformationProvider{}
}

func (p *BotInformationProvider) Version() string {
	return Version
}

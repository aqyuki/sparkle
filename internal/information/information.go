package information

const defaultVersion = "unknown"

// Version is the version of the application.
// This is set at build time using the -ldflags "-X 'internal/info.Version=$VERSION'"
var Version = defaultVersion

// BotInformation is a struct to provide bot information.
type BotInformation struct {
	Version string
}

func NewBotInformation() *BotInformation {
	return &BotInformation{
		Version: Version,
	}
}

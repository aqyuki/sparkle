package bot

import "github.com/bwmarrin/discordgo"

type ReadyHandler interface {
	Handle(*discordgo.Session, *discordgo.Ready)
}

type MessageCreateHandler interface {
	Handle(*discordgo.Session, *discordgo.MessageCreate)
}

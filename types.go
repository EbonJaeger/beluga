package beluga

import (
	"github.com/bwmarrin/discordgo"
)

// Command represents a chat command for a Beluga plugin to handle
type Command struct {
	ChannelID    string
	Command      string
	GuildID      string
	Message      string
	MessageNoCmd string
	Sender       *discordgo.User
}

// CommandPlugin is the interface that plugins can implement to handle commands
// from a Discord channel
type CommandPlugin interface {
	Handle(s *discordgo.Session, c Command)
}

// UserBlacklist contains the user ID's that are blacklisted from using commands
type UserBlacklist struct {
	// Users is a list of blacklisted user ID's
	Users []string
}

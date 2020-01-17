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

// UserBlacklist contains the user ID's that are blacklisted from using commands
type UserBlacklist struct {
	// Users is a list of blacklisted user ID's
	Users []string
}

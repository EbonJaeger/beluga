package plugins

import (
	"github.com/bwmarrin/discordgo"
)

// BelugaCommand represents a chat command for a Beluga
// plugin to handle
type BelugaCommand struct {
	ChannelID    string
	Command      string
	Message      string
	MessageNoCmd string
	Sender       *discordgo.User
}

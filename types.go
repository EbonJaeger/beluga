package beluga

import (
	"github.com/bwmarrin/discordgo"
)

// Command represents a chat command for a Beluga plugin to handle
type Command struct {
	ChannelID    string
	Command      string
	Message      string
	MessageNoCmd string
	Sender       *discordgo.User
}

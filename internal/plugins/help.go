package plugins

import (
	"fmt"
	"strings"

	"github.com/EbonJaeger/beluga"
	"github.com/EbonJaeger/beluga/config"
	"github.com/bwmarrin/discordgo"
)

type helpPlugin struct{}

// Help is our help responder
var Help helpPlugin

func (h *helpPlugin) Handle(s *discordgo.Session, c beluga.Command) {
	// Ignore other commands
	if c.Command != "help" {
		return
	}

	// Build the help response
	var b strings.Builder
	b.WriteString(" **Commands for Beluga Bot:**\n")
	b.WriteString("> `!help` - Show this help message\n")
	if beluga.ArrayContains(config.Conf.Plugins.Enabled, "Hunter2") {
		b.WriteString("> `!hunter2` - Something secret might happen if you use this command...\n")
	}
	if beluga.ArrayContains(config.Conf.Plugins.Enabled, "Slap") {
		b.WriteString("> `!slap <user>` - Creatively slap a user. Name can be a mention, or part of their username\n")
	}

	// Create a DM channel
	dm, _ := s.UserChannelCreate(c.Sender.ID)
	// Attempt to send the response
	if _, err := s.ChannelMessageSend(dm.ID, b.String()); err != nil {
		s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("%s I can't slide into your DM's, are they open?", c.Sender.Mention()))
	}
}

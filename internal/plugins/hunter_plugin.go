package plugins

import (
	"fmt"
	"time"

	"github.com/EbonJaeger/beluga"
	"github.com/bwmarrin/discordgo"
)

// BelugaHunterPlugin is our hunter2 plugin
type BelugaHunterPlugin struct{}

// HunterPlugin is our hunter2 plugin instance
var HunterPlugin BelugaHunterPlugin

// Handle handles the "!hunter2" command
func (p *BelugaHunterPlugin) Handle(s *discordgo.Session, c beluga.Command) {
	// Check that it's the right command
	if c.Command != "hunter2" {
		return
	}

	// Respond
	s.ChannelMessageSend(c.ChannelID, "Validating code, please stand by...")
	time.Sleep(3 * time.Second)
	s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("%s, code validated! :slight_smile:", c.Sender.Mention()))
}

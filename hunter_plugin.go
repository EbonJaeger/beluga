package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// BelugaHunterPlugin is our hunter2 plugin
type BelugaHunterPlugin struct{}

// HunterPlugin is our hunter2 plugin instance
var HunterPlugin BelugaHunterPlugin

// Handle handles the "!hunter2" command
func (p *BelugaHunterPlugin) Handle(s *discordgo.Session, c Command) {
	// Check that it's the right command
	if c.Command != "hunter2" {
		return
	}

	// Respond
	_, err := s.ChannelMessageSend(c.ChannelID, "Validating code, please stand by...")
	if err != nil {
		Log.Errorf("Error responding to '%s' command: %s\n", c.Command, err.Error())
	}
	time.Sleep(3 * time.Second)
	_, err = s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("%s, code validated! :slight_smile:", c.Sender.Mention()))
	if err != nil {
		Log.Errorf("Error responding to '%s' command: %s\n", c.Command, err.Error())
	}
}

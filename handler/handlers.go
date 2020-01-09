package handler

import (
	log "github.com/DataDrake/waterlog"
	"github.com/EbonJaeger/beluga/plugins"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Funcs holds our handler functions with a reference to the
// Beluga plugin manager
type Funcs struct {
	PluginManager *plugins.PluginManager
}

// OnReady handles the "ready" event from Discord
func (f *Funcs) OnReady(s *discordgo.Session, e *discordgo.Ready) {
	s.UpdateStatus(0, "with bits and bobs")
}

// OnGuildCreate handles when we join a Discord guild
func (f *Funcs) OnGuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	// Make sure the guild is available
	if e.Guild.Unavailable {
		log.Warnln("Attempted to join Guild '%s', but it was unavailable")
		return
	}

	// Join the correct channel
	for _, channel := range e.Guild.Channels {
		if channel.ID == e.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Beluga is ready!")
			return
		}
	}
}

// OnMessageCreate handles when a regular message is sent in a channel
// that we have access to
func (f *Funcs) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by ourselves
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Get the message content
	msg := m.Message.Content
	// Check if the message content has more than one character
	if len(msg) < 2 {
		return
	}
	// Split on whitespace
	parts := strings.Split(msg, " ")
	// Check if the message starts with a command prefix
	if strings.HasPrefix(parts[0], "!") {
		// Get the command word
		cmd := strings.Replace(parts[:1][0], "!", "", -1)
		// Make a BelugaCommand
		var bm = plugins.BelugaCommand{
			ChannelID:    m.Message.ChannelID,
			Command:      cmd,
			Message:      msg,
			MessageNoCmd: strings.TrimPrefix(msg, "!"+cmd),
			Sender:       m.Message.Author,
		}
		// Send the command to all handlers
		f.PluginManager.SendCommand(bm)
	}
}

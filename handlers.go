package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// OnReady handles the "ready" event from Discord
func OnReady(s *discordgo.Session, e *discordgo.Ready) {
	s.UpdateStatus(0, "with bits and bobs")
}

// OnGuildCreate handles when we join a Discord guild
func OnGuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	// Make sure the guild is available
	if e.Guild.Unavailable {
		Log.Warnln("Attempted to join Guild '%s', but it was unavailable")
		return
	}

	Log.Infof("Connected to the '%s' guild\n", e.Guild.Name)

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
func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by ourselves
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore message sent from other bots (fun as that would be...)
	if m.Author.Bot {
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
		// Make a Command
		var bm = Command{
			ChannelID:    m.Message.ChannelID,
			Command:      strings.ToLower(cmd),
			Message:      msg,
			MessageNoCmd: strings.TrimPrefix(msg, "!"+cmd),
			Sender:       m.Message.Author,
		}
		// Send the command to all handlers
		PluginManager.SendCommand(bm)
	}
}

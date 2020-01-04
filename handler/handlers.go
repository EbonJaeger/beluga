package handler

import (
	log "github.com/DataDrake/waterlog"
	"github.com/bwmarrin/discordgo"
)

// OnReady handles the "ready" event from Discord
func OnReady(s *discordgo.Session, e *discordgo.Ready) {
	s.UpdateStatus(0, "with bits and bobs")
}

// OnGuildCreate handles when we join a Discord guild
func OnGuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
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
func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by ourselves
	if m.Author.ID == s.State.User.ID {
		return
	}

	// TODO: Implement actual message handling
	log.Infof("Message recieved in channel '%s': %s\n", m.ChannelID, m.Content)
}

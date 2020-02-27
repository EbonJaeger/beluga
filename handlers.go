package beluga

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isDirectMessage(s *discordgo.Session, channelID string) bool {
	// Get the channel
	c, err := s.State.Channel(channelID)
	if err != nil {
		if c, err = s.Channel(channelID); err != nil {
			return false
		}
	}
	// Check the channel type
	return c.Type == discordgo.ChannelTypeDM || c.Type == discordgo.ChannelTypeGroupDM
}

// OnReady handles the "ready" event from Discord
func OnReady(s *discordgo.Session, e *discordgo.Ready) {
	s.UpdateStatus(0, "!help to list commands")
}

// OnGuildCreate handles when we join a Discord guild
func OnGuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	// Make sure the guild is available
	if e.Guild.Unavailable {
		Log.Warnf("Attempted to join {%s}, but it was unavailable\n", e.Guild.ID)
		return
	}

	// Check if this Guild already has a configuration
	if _, exists := Config.Guilds[e.Guild.ID]; !exists {
		// Set defaults for this guild if no configuration is found
		Log.Infof("Creating default configuration for {%s}\n", e.Guild.ID)
		SetGuildDefaults(e.Guild.ID)
	}

	Log.Infof("Joined Guild: {%s}\n", e.Guild.ID)

	cID := e.Guild.SystemChannelID
	s.ChannelMessageSend(cID, "Beluga is ready! Type `!help` for commands list!")
}

// OnMessageCreate handles when a regular message is sent in a channel
// that we have access to
func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from DM's
	if isDirectMessage(s, m.ChannelID) {
		return
	}

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
		// Check if the user is currently blacklisted
		if ArrayContains(Blacklist.Guilds[m.Message.GuildID], m.Author.ID) {
			return
		}
		// Get the command word
		cmd := strings.Replace(parts[:1][0], "!", "", -1)
		// Trim trailing whitespace
		msg = strings.TrimSpace(msg)
		msgNoCmd := strings.TrimSpace(strings.TrimPrefix(msg, "!"+cmd))
		// Make a Command
		var bm = Command{
			ChannelID:    m.Message.ChannelID,
			Command:      strings.ToLower(cmd),
			GuildID:      m.Message.GuildID,
			Message:      msg,
			MessageNoCmd: msgNoCmd,
			Sender:       m.Message.Author,
		}
		// Send the command to all handlers
		Manager.SendCommand(bm)
	}
}

package beluga

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type helpPlugin struct{}

// Help is our help responder
var Help helpPlugin

func (h *helpPlugin) Handle(s *discordgo.Session, c Command) {
	// Ignore other commands
	if c.Command != "help" {
		return
	}

	// Build the help response
	var b strings.Builder
	b.WriteString(" **Commands for Beluga Bot:**\n")
	b.WriteString("> `!help` - Show this help message\n")
	if ArrayContains(Conf.Guilds[c.GuildID].EnabledPlugins, "Hunter2") {
		b.WriteString("> `!hunter2` - Something secret might happen if you use this command...\n")
	}
	if ArrayContains(Conf.Guilds[c.GuildID].EnabledPlugins, "Slap") {
		b.WriteString("> `!slap <user>` - Creatively slap a user. Name can be a mention, or part of their username\n")
	}
	if MemberHasPermission(s, c.GuildID, c.Sender.ID, AdministratorPerm) {
		b.WriteString("**Administrator Commands:**\n")
		b.WriteString("> `!blacklist <user>` - Add a user to the blacklist so they can't run any bot commands\n")
		b.WriteString("> `!rmblacklist <user>` - Remove a user from the blacklist\n")
		b.WriteString("> `!enableplugin <plugin>` - Enable a plugin for this server")
		b.WriteString("> `!disableplugin <plugin>` - Disable a plugin for this server")
		if ArrayContains(Conf.Guilds[c.GuildID].EnabledPlugins, "Commands") {
			b.WriteString("> `!addcommand <command> <response>` - Add a custom command that the bot will respond to\n")
			b.WriteString("> `!rmcommand <command>` - Remove a custom command\n")
		}
	}

	// Create a DM channel
	dm, _ := s.UserChannelCreate(c.Sender.ID)
	// Attempt to send the response
	if _, err := s.ChannelMessageSend(dm.ID, b.String()); err != nil {
		s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("%s I can't slide into your DM's, are they open?", c.Sender.Mention()))
	}
}

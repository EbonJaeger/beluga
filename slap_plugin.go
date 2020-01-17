package beluga

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// SlapPlugin is our slapping plugin
type SlapPlugin struct{}

// Slapper is our slap plugin instance
var Slapper SlapPlugin

// Handle implements the '!slap' command
func (p *SlapPlugin) Handle(s *discordgo.Session, c Command) {
	// Ignore other commands
	if c.Command != "slap" {
		return
	}
	// Check for args
	if len(c.MessageNoCmd) > 0 {
		// Split args
		args := strings.Split(c.MessageNoCmd, " ")
		// Check args length
		if len(args) == 1 {
			// Get target user's name
			name := args[0]
			// Get the Guild
			g, _ := s.Guild(c.GuildID)
			// Find the target user
			target := GetUserFromName(s, g, name)
			// Make sure we have a valid target
			if target == nil {
				s.ChannelMessageSend(c.ChannelID, "You must be halucinating. There is noone here by that name.")
				return
			}
			var resp string
			// Check for self-harm
			if target.Username == c.Sender.Username {
				s.ChannelMessageSend(c.ChannelID, Conf.Guilds[c.GuildID].SlapConfig.SelfSlap)
				return
			}
			// Check for slapping Beluga bot
			if target.Username == "Beluga" {
				target = c.Sender
			}
			// Seed random
			rand.Seed(time.Now().Unix())
			// Get random number
			randNum := rand.Intn(len(Conf.Guilds[c.GuildID].SlapConfig.SlapMessages))
			// Get our slap message
			resp = "*" + Conf.Guilds[c.GuildID].SlapConfig.SlapMessages[randNum] + "*"
			// Put in the target's @-mention
			resp = strings.Replace(resp, "$USER", target.Mention(), -1)
			// Send the response
			s.ChannelMessageSend(c.ChannelID, resp)
		}
	}
}

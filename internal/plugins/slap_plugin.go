package plugins

import (
	"math/rand"
	"strings"
	"time"

	"github.com/EbonJaeger/beluga"
	"github.com/EbonJaeger/beluga/config"
	"github.com/bwmarrin/discordgo"
)

// SlapPlugin is our slapping plugin
type SlapPlugin struct{}

// Slapper is our slap plugin instance
var Slapper SlapPlugin

// Handle implements the '!slap' command
func (p *SlapPlugin) Handle(s *discordgo.Session, c beluga.Command) {
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
			// Get target user's ID
			name := args[0]
			// Get the Guild
			g, _ := s.Guild(c.GuildID)
			// Find the target user
			target := parseTarget(s, g, name)
			// Make sure we have a valid target
			if target == nil {
				s.ChannelMessageSend(c.ChannelID, "You must be halucinating. There is noone here by that name.")
				return
			}
			var resp string
			// Check for self-harm
			if target.Username == c.Sender.Username {
				s.ChannelMessageSend(c.ChannelID, config.Conf.SlapConf.SelfSlap)
				return
			}
			// Check for slapping Beluga bot
			if target.Username == "Beluga" {
				target = c.Sender
			}
			// Seed random
			rand.Seed(time.Now().Unix())
			// Get random number
			randNum := rand.Intn(len(config.Conf.SlapConf.SlapMessages))
			// Get our slap message
			resp = "*" + config.Conf.SlapConf.SlapMessages[randNum] + "*"
			// Put in the target's @-mention
			resp = strings.Replace(resp, "$USER", target.Mention(), -1)
			// Send the response
			s.ChannelMessageSend(c.ChannelID, resp)
		}
	}
}

func parseTarget(s *discordgo.Session, g *discordgo.Guild, t string) *discordgo.User {
	var target *discordgo.User
	// Check if it's a mention
	if strings.HasPrefix(t, "<@!") {
		id := strings.TrimPrefix(t, "<@!")
		id = strings.TrimSuffix(id, ">")
		target, _ = s.User(id)
	} else {
		// Check if it's a partial name
		for _, u := range g.Members {
			if strings.Contains(strings.ToLower(u.User.Username), strings.ToLower(t)) {
				target = u.User
				break
			}
		}
	}

	return target
}

package beluga

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AdminPlugin is our admin plugin struct
type AdminPlugin struct{}

// BelugaAdmin is our admin plugin instance
var BelugaAdmin AdminPlugin

// Handle handles all admin-related commands
func (p *AdminPlugin) Handle(s *discordgo.Session, c Command) {
	// Send to the right sub-handler
	switch c.Command {
	case "blacklist":
		addBlacklistedUser(s, c)
		break
	case "rmblacklist":
		removeBlacklistedUser(s, c)
		break
	default:
		return
	}
}

func addBlacklistedUser(s *discordgo.Session, c Command) {
	// Check if the sender is an administrator
	if !MemberHasPermission(s, c.GuildID, c.Sender.ID, AdministratorPerm) {
		s.ChannelMessageSend(c.ChannelID, "You don't have permission to perform that command! Get outa here! :angry:")
		return
	}
	// Check for args
	if len(c.MessageNoCmd) > 0 {
		// Split args
		args := strings.Split(c.MessageNoCmd, " ")
		// Check args length
		if len(args) == 1 {
			raw := args[0]
			// Get the channel's Guild
			guild, _ := s.Guild(c.GuildID)
			// Get the target user
			user := GetUserFromName(s, guild, raw)
			// Make sure we have a valid target
			if user == nil {
				s.ChannelMessageSend(c.ChannelID, "You must be halucinating. There is noone here by that name.")
				return
			}
			// Check if the user is an admin
			if MemberHasPermission(s, c.GuildID, user.ID, AdministratorPerm) {
				s.ChannelMessageSend(c.ChannelID, "I cant blacklist that user, are you kidding? :frowning:")
				return
			}
			// Check if user is already blacklisted
			if ArrayContains(Blacklist.Users, user.ID) {
				s.ChannelMessageSend(c.ChannelID, "That user is already blacklisted!")
				return
			}
			// Add user to blacklist
			Blacklist.Users = append(Blacklist.Users, user.ID)
			SaveBlacklist()
		}
	}
}

func removeBlacklistedUser(s *discordgo.Session, c Command) {
	// Check if the sender is an administrator
	if MemberHasPermission(s, c.GuildID, c.Sender.ID, AdministratorPerm) {
		s.ChannelMessageSend(c.ChannelID, "You don't have permission to perform that command! Get outa here! :angry:")
		return
	}
	// Check for args
	if len(c.MessageNoCmd) > 0 {
		// Split args
		args := strings.Split(c.MessageNoCmd, " ")
		// Check args length
		if len(args) == 1 {
			raw := args[0]
			// Get the channel's Guild
			guild, _ := s.Guild(c.GuildID)
			// Get the target user
			user := GetUserFromName(s, guild, raw)
			// Make sure we have a valid target
			if user == nil {
				s.ChannelMessageSend(c.ChannelID, "You must be halucinating. There is noone here by that name.")
				return
			}
			// Check if user is actually blacklisted
			if !ArrayContains(Blacklist.Users, user.ID) {
				s.ChannelMessageSend(c.ChannelID, "That user isn't blacklisted!")
				return
			}
			// Remove user from blacklist
			Blacklist.Users = RemoveFromStringArray(Blacklist.Users, user.ID)
			SaveBlacklist()
		}
	}
}

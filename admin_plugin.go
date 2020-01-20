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
	// Check if the sender is an administrator
	if !MemberHasPermission(s, c.GuildID, c.Sender.ID, AdministratorPerm) {
		s.ChannelMessageSend(c.ChannelID, "You don't have permission to perform that command! Get outa here! :angry:")
		return
	}
	// Send to the right sub-handler
	switch c.Command {
	case "blacklist":
		addBlacklistedUser(s, c)
		break
	case "disableplugin":
		disablePlugin(s, c)
		break
	case "enableplugin":
		enablePlugin(s, c)
		break
	case "rmblacklist":
		removeBlacklistedUser(s, c)
		break
	default:
		return
	}
}

func addBlacklistedUser(s *discordgo.Session, c Command) {
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
			if err := SaveConfigToFile("blacklist.toml", Blacklist); err != nil {
				Log.Errorf("Error while saving blacklist file: %s\n", err.Error())
				s.ChannelMessageSend(c.ChannelID, "An error occurred while saving the blacklist. :frowning:")
			} else {
				s.ChannelMessageSend(c.ChannelID, "Added the user to the blacklist!")
			}
		}
	}
}

func disablePlugin(s *discordgo.Session, c Command) {
	// Check for args
	if len(c.MessageNoCmd) > 0 {
		// Split args
		args := strings.Split(c.MessageNoCmd, " ")
		// Check args length
		if len(args) == 1 {
			raw := args[0]
			// Get the channel's Guild
			guild, _ := s.Guild(c.GuildID)
			// Capitalize the first letter
			plugin := strings.Title(raw)
			// Check if the plugin exists and is enabled
			if Manager.IsLoaded(plugin) && Manager.IsEnabled(guild.ID, plugin) {
				// Remove the plugin from the guild config
				Conf.Guilds[guild.ID].EnabledPlugins = RemoveFromStringArray(Conf.Guilds[guild.ID].EnabledPlugins, plugin)
				// Save config to file
				if err := SaveConfigToFile("beluga.conf", Conf); err == nil {
					s.ChannelMessageSend(c.ChannelID, "Plugin disabled! :smiley:")
				} else {
					Log.Errorf("Error saving config file: %s\n", err.Error())
					s.ChannelMessageSend(c.ChannelID, "There was a problem saving the config. :frowning:")
				}
			} else {
				s.ChannelMessageSend(c.ChannelID, "I don't know what that plugin is! :frowning:")
			}
		}
	}
}

func enablePlugin(s *discordgo.Session, c Command) {
	// Check for args
	if len(c.MessageNoCmd) > 0 {
		// Split args
		args := strings.Split(c.MessageNoCmd, " ")
		// Check args length
		if len(args) == 1 {
			raw := args[0]
			// Get the channel's Guild
			guild, _ := s.Guild(c.GuildID)
			// Capitalize the first letter
			plugin := strings.Title(raw)
			// Check if the plugin exists and isn't already enabled
			if Manager.IsLoaded(plugin) && !Manager.IsEnabled(c.GuildID, plugin) {
				// Add the plugin to the guild config
				Conf.Guilds[guild.ID].EnabledPlugins = append(Conf.Guilds[guild.ID].EnabledPlugins, plugin)
				// Save config to file
				if err := SaveConfigToFile("beluga.conf", Conf); err == nil {
					s.ChannelMessageSend(c.ChannelID, "Plugin enabled! :smiley:")
				} else {
					Log.Errorf("Error saving config file: %s\n", err.Error())
					s.ChannelMessageSend(c.ChannelID, "There was a problem saving the config. :frowning:")
				}
			} else {
				s.ChannelMessageSend(c.ChannelID, "I don't know what that plugin is! :frowning:")
			}
		}
	}
}

func removeBlacklistedUser(s *discordgo.Session, c Command) {
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
			if err := SaveConfigToFile("blacklist.toml", Blacklist); err != nil {
				Log.Errorf("Error while saving blacklist file: %s\n", err.Error())
				s.ChannelMessageSend(c.ChannelID, "An error occurred while saving the blacklist. :frowning:")
			} else {
				s.ChannelMessageSend(c.ChannelID, "Removed the user to the blacklist! :smiley:")
			}
		}
	}
}

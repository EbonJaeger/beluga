package beluga

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// AdminPlugin is our admin plugin struct
type AdminPlugin struct {
	blMu sync.Mutex // Blacklist lock
	cMu  sync.Mutex // Config lock
}

// BelugaAdmin is our admin plugin instance
var BelugaAdmin AdminPlugin

// Handle handles all admin-related commands
func (p *AdminPlugin) Handle(s *discordgo.Session, c Command) {
	// Check if the sender is an administrator
	if !MemberHasPermission(s, c.GuildID, c.Sender.ID, discordgo.PermissionAdministrator) {
		s.ChannelMessageSend(c.ChannelID, "You don't have permission to perform that command! Get outa here! :angry:")
		return
	}
	// Send to the right sub-handler
	switch c.Command {
	case "blacklist":
		p.addBlacklistedUser(s, c)
		break
	case "disableplugin":
		p.disablePlugin(s, c)
		break
	case "enableplugin":
		p.enablePlugin(s, c)
		break
	case "listplugins":
		p.listPlugins(s, c)
		break
	case "rmblacklist":
		p.removeBlacklistedUser(s, c)
		break
	default:
		return
	}
}

func (p *AdminPlugin) addBlacklistedUser(s *discordgo.Session, c Command) {
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
			if MemberHasPermission(s, c.GuildID, user.ID, discordgo.PermissionAdministrator) {
				s.ChannelMessageSend(c.ChannelID, "I cant blacklist that user, are you kidding? :frowning:")
				return
			}
			// Grab a lock on the blacklist
			p.blMu.Lock()
			defer p.blMu.Unlock()
			// Check if user is already blacklisted
			if ArrayContains(Blacklist.Guilds[guild.ID], user.ID) {
				s.ChannelMessageSend(c.ChannelID, "That user is already blacklisted!")
				return
			}
			// Add user to blacklist
			Blacklist.Guilds[guild.ID] = append(Blacklist.Guilds[guild.ID], user.ID)
			if err := SaveConfigToFile("blacklist.toml", Blacklist); err != nil {
				Log.Errorf("Error while saving blacklist file: %s\n", err.Error())
				s.ChannelMessageSend(c.ChannelID, "An error occurred while saving the blacklist. :frowning:")
			} else {
				s.ChannelMessageSend(c.ChannelID, "Added the user to the blacklist!")
			}
		}
	}
}

func (p *AdminPlugin) disablePlugin(s *discordgo.Session, c Command) {
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
				// Grab a lock on the config
				p.cMu.Lock()
				defer p.cMu.Unlock()
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

func (p *AdminPlugin) enablePlugin(s *discordgo.Session, c Command) {
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
				// Grab a lock on the config
				p.cMu.Lock()
				defer p.cMu.Unlock()
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

func (p *AdminPlugin) listPlugins(s *discordgo.Session, c Command) {
	// Grab a lock on the config
	p.cMu.Lock()
	// Get the enabled plugins for this guild
	e := Conf.Guilds[c.GuildID].EnabledPlugins
	// Get all available plugins
	t := make([]string, len(Manager.Plugins))
	i := 0
	for k := range Manager.Plugins {
		t[i] = k
		i++
	}
	// Unlock again
	p.cMu.Unlock()
	// Remove "always-on" plugins like help and admin
	t = RemoveMultipleFromArray(t, []string{"Admin", "Help"})
	// Remove enabled plugins from the total list
	d := RemoveMultipleFromArray(t, e)
	// Join the lists into strings
	enabled := strings.Join(e, ", ")
	disabled := strings.Join(d, ", ")
	// Get the name of the Guild
	guild, _ := s.Guild(c.GuildID)
	// Create a single string response
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Plugins for %s:\n", guild.Name))
	b.WriteString(fmt.Sprintf("> Enabled: %s\n", enabled))
	b.WriteString(fmt.Sprintf("> Disabled: %s", disabled))
	resp := b.String()
	// Create a DM with the sender
	dm, _ := s.UserChannelCreate(c.Sender.ID)
	// DM the response
	s.ChannelMessageSend(dm.ID, resp)
}

func (p *AdminPlugin) removeBlacklistedUser(s *discordgo.Session, c Command) {
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
			// Grab a lock on the blacklist
			p.blMu.Lock()
			defer p.blMu.Unlock()
			// Check if user is actually blacklisted
			if !ArrayContains(Blacklist.Guilds[guild.ID], user.ID) {
				s.ChannelMessageSend(c.ChannelID, "That user isn't blacklisted!")
				return
			}
			// Remove user from blacklist
			Blacklist.Guilds[guild.ID] = RemoveFromStringArray(Blacklist.Guilds[guild.ID], user.ID)
			if err := SaveConfigToFile("blacklist.toml", Blacklist); err != nil {
				Log.Errorf("Error while saving blacklist file: %s\n", err.Error())
				s.ChannelMessageSend(c.ChannelID, "An error occurred while saving the blacklist. :frowning:")
			} else {
				s.ChannelMessageSend(c.ChannelID, "Removed the user to the blacklist! :smiley:")
			}
		}
	}
}

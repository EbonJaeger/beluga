package beluga

import "github.com/bwmarrin/discordgo"

import "strings"

import "fmt"

// CustomCommandsPlugin is our custom command creater and responder plugin
type CustomCommandsPlugin struct{}

// Commands is our ResponderPlugin instance
var Commands CustomCommandsPlugin

// Handle handles the commands for the responder plugin
func (p *CustomCommandsPlugin) Handle(s *discordgo.Session, c Command) {
	switch c.Command {
	case "addcommand":
		if MemberHasPermission(s, c.GuildID, c.Sender.ID, discordgo.PermissionAdministrator) {
			if len(c.MessageNoCmd) > 0 {
				// Split out the first word
				parts := strings.SplitN(c.MessageNoCmd, " ", 2)
				cmd := parts[0]
				resp := parts[1]
				// Make sure we have a map to add to
				if Config.Guilds[c.GuildID].CustomResponses == nil {
					Config.Guilds[c.GuildID].CustomResponses = make(map[string]string)
				}
				// Add the command to the config
				Config.Guilds[c.GuildID].CustomResponses[cmd] = resp
				// Save the config to file
				if err := SaveConfigToFile("beluga.conf", Config); err != nil {
					Log.Errorf("Error while saving config: %s\n", err.Error())
					s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("Error while updating command '%s' :frowning:", cmd))
				} else {
					s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("Command '%s' added or updated successfully! :smiley:", cmd))
				}
			}
		} else {
			s.ChannelMessageSend(c.ChannelID, "You don't have permission to perform that command! Get outa here! :angry:")
		}
		break
	case "rmcommand":
		if MemberHasPermission(s, c.GuildID, c.Sender.ID, discordgo.PermissionAdministrator) {
			if len(c.MessageNoCmd) > 0 {
				// Get the command
				cmd := strings.SplitN(c.MessageNoCmd, " ", 1)[0]
				if cmd != "" {
					// Remove the key from the responses map
					delete(Config.Guilds[c.GuildID].CustomResponses, cmd)
					// Save the config to file
					if err := SaveConfigToFile("beluga.conf", Config); err != nil {
						Log.Errorf("Error while saving config: %s\n", err.Error())
						s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("Error while removing command '%s' :frowning:", cmd))
					} else {
						s.ChannelMessageSend(c.ChannelID, fmt.Sprintf("Command '%s' removed successfully! :smiley:", cmd))
					}
				}
			}
		} else {
			s.ChannelMessageSend(c.ChannelID, "You don't have permission to perform that command! Get outa here! :angry:")
		}
		break
	default:
		// Get all custom commands for the current Guild
		commands := Config.Guilds[c.GuildID].CustomResponses
		if resp := commands[c.Command]; resp != "" {
			s.ChannelMessageSend(c.ChannelID, resp)
		}
	}
}

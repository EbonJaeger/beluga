package beluga

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	// PluginsPath is the path where Beluga plugins should
	// be placed
	PluginsPath = "/usr/share/beluga/plugins"
)

// PluginManager is the Beluga plugin manager
type PluginManager struct {
	Plugins map[string]plugin.Symbol
}

// IsEnabled will check if the given plugin is enabled in the
// Beluga config
func (pm *PluginManager) IsEnabled(guild string, name string) bool {
	return ArrayContains(Conf.Guilds[guild].EnabledPlugins, name)
}

// LoadPlugins attempts to load all found plugins
func (pm *PluginManager) LoadPlugins() error {
	var pluginLoadErr error

	Log.Infoln("Looking for plugins to enable")

	// Open plugin directory
	if pluginDir, err := os.Open(PluginsPath); err == nil {
		defer pluginDir.Close()
		// Read directory contents
		if children, readErr := pluginDir.Readdir(-1); readErr == nil {
			// Check for contents
			if len(children) > 0 {
				for _, child := range children {
					// Get file name and extension
					fileName := child.Name()
					fileExt := filepath.Ext(fileName)

					// Make sure it's a library file (.so)
					if !child.IsDir() && (fileExt == ".so") {
						// Get the plugin name
						pluginName := strings.Replace(fileName, fileExt, "", -1)
						// Make sure we haven't already added this plugin
						if _, added := pm.Plugins[pluginName]; !added {
							// Open the file
							if plugin, err := plugin.Open(filepath.Join(PluginsPath, fileName)); err == nil {
								Log.Infof("Checking '%s' for a message handler\n", fileName)
								// Look for message handler function
								if handleFunc, err := plugin.Lookup("Handle"); err == nil {
									Log.Goodf("Added plugin '%s'\n", pluginName)
									// Add the plugin
									pm.Plugins[pluginName] = handleFunc
								} else {
									Log.Warnf("Error while loading plugin '%s': %s\n", pluginName, err.Error())
									continue
								}
							} else {
								Log.Warnf("Error while loading plugin '%s': %s\n", pluginName, err.Error())
								continue
							}
						}
					}
				}
			} else {
				Log.Infoln("No plugins found")
			}
		} else {
			pluginLoadErr = readErr
		}
	} else {
		pluginLoadErr = err
	}
	return pluginLoadErr
}

// SendCommand sends a chat command to all registered handlers
func (pm *PluginManager) SendCommand(cmd Command) {
	// Send to help handler
	Help.Handle(Session, cmd)

	// Send to admin handler
	BelugaAdmin.Handle(Session, cmd)

	// Send to hunter2 plugin
	if pm.IsEnabled(cmd.GuildID, "Hunter2") {
		Hunter.Handle(Session, cmd)
	}

	// Send to slap plugin
	if pm.IsEnabled(cmd.GuildID, "Slap") {
		Slapper.Handle(Session, cmd)
	}

	// Send to all third-party plugins
	for _, handleFunc := range pm.Plugins {
		handleFunc.(func(*discordgo.Session, Command))(Session, cmd)
	}
}

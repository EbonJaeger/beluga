package plugins

import (
	log "github.com/DataDrake/waterlog"
	"github.com/EbonJaeger/beluga/config"
	"github.com/bwmarrin/discordgo"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

const (
	// PluginsPath is the path where Beluga plugins should
	// be placed
	PluginsPath = "/usr/share/beluga/plugins"
)

// PluginManager is the Beluga plugin manager
type PluginManager struct {
	Plugins map[string]plugin.Symbol
	Session *discordgo.Session
}

// NewManager creates a new plugin manager
func NewManager() *PluginManager {
	return &PluginManager{}
}

// IsEnabled will check if the given plugin is enabled in the
// Beluga config
func (pm *PluginManager) IsEnabled(name string) bool {
	return ArrayContains(config.Conf.Plugins, name)
}

// LoadPlugins attempts to load all found plugins
func (pm *PluginManager) LoadPlugins() error {
	var pluginLoadErr error

	log.Infoln("Looking for plugins to enable")

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
							if plugin, openErr := plugin.Open(filepath.Join(PluginsPath, fileName)); openErr == nil {
								log.Infof("Checking '%s' for a message handler\n", fileName)
								// Look for message handler function
								if handleFunc, lookupErr := plugin.Lookup("Handle"); lookupErr == nil {
									log.Goodf("Added plugin '%s'\n", pluginName)
									// Add the plugin
									pm.Plugins[pluginName] = handleFunc
								} else {
									pluginLoadErr = lookupErr
									break
								}
							} else {
								pluginLoadErr = openErr
								break
							}
						}
					}
				}
			} else {
				log.Infoln("No plugins found")
			}
		} else {
			pluginLoadErr = readErr
		}

	}
	return pluginLoadErr
}

// SendCommand sends a chat command to all registered handlers
func (pm *PluginManager) SendCommand(cmd BelugaCommand) {
	for _, handleFunc := range pm.Plugins {
		handleFunc.(func(*discordgo.Session, BelugaCommand))(pm.Session, cmd)
	}
}

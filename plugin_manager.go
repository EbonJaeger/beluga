package main

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/EbonJaeger/beluga/config"
	"github.com/bwmarrin/discordgo"
)

const (
	// PluginsPath is the path where Beluga plugins should
	// be placed
	PluginsPath = "/usr/share/beluga/plugins"
)

// BelugaPluginManager is the Beluga plugin manager
type BelugaPluginManager struct {
	Plugins map[string]plugin.Symbol
}

// IsEnabled will check if the given plugin is enabled in the
// Beluga config
func (pm *BelugaPluginManager) IsEnabled(name string) bool {
	return ArrayContains(config.Conf.Plugins, name)
}

// LoadPlugins attempts to load all found plugins
func (pm *BelugaPluginManager) LoadPlugins() error {
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
							if plugin, openErr := plugin.Open(filepath.Join(PluginsPath, fileName)); openErr == nil {
								Log.Infof("Checking '%s' for a message handler\n", fileName)
								// Look for message handler function
								if handleFunc, lookupErr := plugin.Lookup("Handle"); lookupErr == nil {
									Log.Goodf("Added plugin '%s'\n", pluginName)
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
func (pm *BelugaPluginManager) SendCommand(cmd Command) {
	// Send to hunter2 plugin
	if PluginManager.IsEnabled("Hunter2") {
		HunterPlugin.Handle(Session, cmd)
	}

	// Send to all third-party plugins
	for _, handleFunc := range pm.Plugins {
		handleFunc.(func(*discordgo.Session, Command))(Session, cmd)
	}
}

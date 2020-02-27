package beluga

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// RootConfig is our root configuration structure
type RootConfig struct {
	// Token is the Discord bot token to use to connect to Discord
	Token string `toml:"discord_bot_token"`
	// Guilds contains the guild-specific configurations for each guild
	Guilds map[string]*GuildConfig
}

// GuildConfig is a guild-specific configuration
type GuildConfig struct {
	EnabledPlugins  []string
	CustomResponses map[string]string
	SlapConfig      SlapPluginConfig
}

// SlapPluginConfig holds the configuration options for the slap plugin.
// These options are guild-specific
type SlapPluginConfig struct {
	SelfSlap     string   `toml:"self-slap"`
	SlapMessages []string `toml:"slap-messages,omitempty"`
}

// LoadConfig will load the config and return it
func LoadConfig() (conf RootConfig, err error) {
	// Get our config file
	path := filepath.Join(ConfigDir, "beluga.conf")
	Log.Infof("Using configuration at '%s'\n", path)
	if err = CreateFileIfNotExists(path); err != nil {
		return
	}
	// Parse file
	if _, err = toml.DecodeFile(path, &conf); err != nil {
		// Validate Guild configuration section
		if conf.Guilds == nil {
			conf.Guilds = make(map[string]*GuildConfig)
		}
		return
	}
	return
}

// SetDefaults sets any sane default configuration options
func SetDefaults() RootConfig {
	return RootConfig{
		Token: "",
	}
}

// SetGuildDefaults creates a new guild configuration with defaults, adds it to
// the existing configuration, and saves it to the disk
func SetGuildDefaults(guildID string) {
	// Create a fresh Guild config
	g := GuildConfig{
		EnabledPlugins:  []string{"Hunter2", "Slap"},
		CustomResponses: make(map[string]string),
		SlapConfig: SlapPluginConfig{
			SelfSlap: "I shall not listen to the demands of mere humans, for I am the underwater whale master.",
			SlapMessages: []string{
				"Annihilates $USER",
				"Claps $USER",
				"Decimates $USER",
				"Destroys $USER",
				"Discombobulates $USER",
				"Does far worse, taking $USER's system and (re-)installing Windows",
				"Gives $USER a splinter",
				"Just looks at $USER with disappointment",
				"Opts to not slap $USER today, but rather gives them a cookie",
				"Punches $USER",
				"Slaps $USER",
				"Snaps its flippers together, $USER turns into ash and disappears into the wind",
				"Thinks $USER should lose a few pounds",
				"Throws $USER down a ravine",
			},
		},
	}
	// Add the new guild to the rest of the config
	Config.Guilds[guildID] = &g
	// Save it to the disk
	SaveConfigToFile("beluga.conf", Config)
}

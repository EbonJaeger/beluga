package beluga

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	// Token is the Discord bot token to use to connect to Discord
	Token string `toml:"discord_bot_token"`
	// Guilds contains the guild-specific configurations for each guild
	Guilds map[string]GuildConfig
	Facts  []string `toml:"facts,omitempty"`
}

// GuildConfig is a guild-specific configuration
type GuildConfig struct {
	EnabledPlugins []string
	SlapConfig     SlapPluginConfig
}

// SlapPluginConfig holds the configuration options for the slap plugin.
// These options are guild-specific
type SlapPluginConfig struct {
	SelfSlap     string   `toml:"self-slap"`
	SlapMessages []string `toml:"slap-messages,omitempty"`
}

// Conf is the current configuration
var Conf config

// LoadConfig will load the config
func LoadConfig() error {
	// Get our config file
	path := filepath.Join(ConfigPath, "beluga.conf")
	if err := CreateFileIfNotExists(path); err != nil {
		return err
	}
	// Parse file
	Conf = config{}
	_, err := toml.DecodeFile(path, &Conf)
	if err != nil {
		return err
	}
	// Validate Guild configuration section
	if Conf.Guilds == nil {
		Conf.Guilds = make(map[string]GuildConfig)
	}

	return nil
}

// SetGuildDefaults creates a new guild configuration with defaults, adds it to
// the existing configuration, and saves it to the disk
func SetGuildDefaults(guildID string) {
	// Create a fresh Guild config
	g := GuildConfig{
		EnabledPlugins: []string{"Hunter2", "Slap"},
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
	Conf.Guilds[guildID] = g
	// Save it to the disk
	SaveConfigToFile("beluga.conf", Conf)
}

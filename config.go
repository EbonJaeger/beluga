package beluga

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
)

type config struct {
	Token    string `toml:"discord_bot_token"`
	Plugins  pluginConf
	SlapConf slapConf `toml:"slap-plugin"`
	Facts    []string `toml:"facts,omitempty"`
}

type pluginConf struct {
	Enabled []string `toml:"enabled-plugins"`
}

type slapConf struct {
	SelfSlap     string   `toml:"self-slap"`
	SlapMessages []string `toml:"slap-messages"`
}

// Conf is the current configuration
var Conf config

// LoadConfig will load the config
func LoadConfig() error {
	// Get our config file
	path := filepath.Join(ConfigPath, "beluga.conf")

	// Parse file
	Conf = config{}
	_, err := toml.DecodeFile(path, &Conf)
	if err != nil {
		return err
	}

	return nil
}
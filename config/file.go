package config

import (
	"github.com/BurntSushi/toml"
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

// Load will load the config
func Load() error {
	// Get our config file
	path := "/etc/beluga/beluga.conf"

	// Parse file
	Conf = config{}
	_, err := toml.DecodeFile(path, &Conf)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	if err := Load(); err != nil {
		panic(err.Error())
	}
}

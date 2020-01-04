package config

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Token string `toml:"discord_bot_token"`
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

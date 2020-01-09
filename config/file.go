package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultSocket is the default socket path for the daemon
	DefaultSocket = "/run/beluga.socket"
)

type config struct {
	Socket  string   `toml:"socket"`
	Token   string   `toml:"discord_bot_token"`
	Plugins []string `toml:"plugins,omitempty"`
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

	// Validate socket file
	if len(Conf.Socket) == 0 {
		Conf.Socket = DefaultSocket
		log.Printf("No socket specified. Using default: %s\n", DefaultSocket)
	}

	return nil
}

func init() {
	if err := Load(); err != nil {
		panic(err.Error())
	}
}

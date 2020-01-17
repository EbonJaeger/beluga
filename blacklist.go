package beluga

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ReadBlacklist reads the blacklist file
func ReadBlacklist() (UserBlacklist, error) {
	var blacklist UserBlacklist
	path := filepath.Join(ConfigPath, "blacklist.toml")
	if err := CreateFileIfNotExists(path); err != nil {
		return blacklist, err
	}
	// Parse the file
	_, err := toml.DecodeFile(path, &blacklist)
	if err != nil {
		return blacklist, err
	}

	return blacklist, nil
}

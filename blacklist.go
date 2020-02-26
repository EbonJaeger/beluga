package beluga

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LoadBlacklist reads the blacklist file
func LoadBlacklist() (blacklist UserBlacklist, err error) {
	path := filepath.Join(ConfigDir, "blacklist.toml")
	if err = CreateFileIfNotExists(path); err != nil {
		return
	}
	// Parse the file
	if _, err = toml.DecodeFile(path, &blacklist); err != nil {
		// Make sure we have a proper guilds section
		if blacklist.Guilds == nil {
			blacklist.Guilds = make(map[string][]string)
		}
		return
	}
	return
}

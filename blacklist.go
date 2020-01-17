package beluga

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ReadBlacklist reads the blacklist file
func ReadBlacklist() (UserBlacklist, error) {
	var blacklist UserBlacklist
	path := filepath.Join(ConfigPath, "blacklist.toml")
	// Check if the file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// Create the file if it doesn't exist
			f, createErr := os.Create(path)
			if createErr != nil {
				return blacklist, createErr
			}
			// Set the file permissions
			if chmodErr := f.Chmod(0644); chmodErr != nil {
				return blacklist, chmodErr
			}
		} else {
			// Other error
			return blacklist, err
		}
	}
	// Parse the file
	_, err := toml.DecodeFile(path, &blacklist)
	if err != nil {
		return blacklist, err
	}

	return blacklist, nil
}

package beluga

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

const (
	// AdministratorPerm is the Discord administrator permission value
	AdministratorPerm = 0x8
)

// ArrayContains checks if a given element is in a string
// array
func ArrayContains(arr []string, element string) bool {
	var found bool

	// Iterate over the array
	for _, ele := range arr {
		// Check if it's the same item
		if ele == element {
			found = true
			break
		}
	}

	return found
}

// CreateFileIfNotExists creates a new empty file at the given path it the file
// does not yet exist
func CreateFileIfNotExists(path string) error {
	// Check if the file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// Create the file if it doesn't exist
			f, createErr := os.Create(path)
			if createErr != nil {
				return createErr
			}
			// Set the file permissions
			if chmodErr := f.Chmod(0644); chmodErr != nil {
				return chmodErr
			}
		} else {
			// Other error
			return err
		}
	}
	return nil
}

// GetUserFromName gets the Discord user from a mention or username. The username
// can be only a partial username
func GetUserFromName(s *discordgo.Session, g *discordgo.Guild, t string) *discordgo.User {
	var target *discordgo.User
	// Check if it's a mention
	if strings.HasPrefix(t, "<@!") {
		// Trim the mention prefix and suffix
		id := strings.TrimPrefix(t, "<@!")
		id = strings.TrimSuffix(id, ">")
		// Get the user
		target, _ = s.User(id)
	} else {
		// Look through all users in the guild
		for _, u := range g.Members {
			// Check if the name matches or is a partial
			if strings.Contains(strings.ToLower(u.User.Username), strings.ToLower(t)) {
				target = u.User
				break
			}
		}
	}

	return target
}

// MemberHasPermission checks if a member of a guild has the desired permission
func MemberHasPermission(s *discordgo.Session, guildID string, userID string, perm int) bool {
	// Get the guild member
	m, err := s.State.Member(guildID, userID)
	if err != nil {
		if m, err = s.GuildMember(guildID, userID); err != nil {
			return false
		}
	}
	// Iterate through all roles to check permissions
	for _, roleID := range m.Roles {
		// Get the role
		role, err := s.State.Role(guildID, roleID)
		// Make sure the role exists
		if err != nil {
			return false
		}
		// Check if the role's permissions contains the sought after permission
		if role.Permissions&perm != 0 {
			return true
		}
	}
	return false
}

// RemoveFromStringArray removes an item from a string array
func RemoveFromStringArray(arr []string, item string) []string {
	// Create a new array
	newArr := []string{}
	// Iterate over the given list
	for _, e := range arr {
		// Add to new list if it doesn't match
		if e != item {
			newArr = append(newArr, e)
		}
	}
	return newArr
}

// SaveConfigToFile saves the given data to the given file name in the local
// config directory
func SaveConfigToFile(name string, data interface{}) {
	var (
		buffer  bytes.Buffer
		saveErr error
	)
	path := filepath.Join(ConfigPath, name)
	// Create our buffer and encoder
	writer := bufio.NewWriter(&buffer)
	encoder := toml.NewEncoder(writer)
	// Encode the struct as toml
	if saveErr = encoder.Encode(data); saveErr == nil {
		// Write to the blacklist file
		saveErr = ioutil.WriteFile(path, buffer.Bytes(), 0644)
	}
	// Log if there's an error
	if saveErr != nil {
		Log.Errorf("Error saving to file: %s\n", saveErr.Error())
	}
}

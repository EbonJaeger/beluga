package beluga

import (
	"strings"

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
func MemberHasPermission(s *discordgo.Session, guildID string, userID string, perm int) (bool, error) {
	// Get the guild member
	m, err := s.State.Member(guildID, userID)
	if err != nil {
		if m, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}
	// Iterate through all roles to check permissions
	for _, roleID := range m.Roles {
		// Get the role
		role, err := s.State.Role(guildID, roleID)
		// Make sure the role exists
		if err != nil {
			return false, err
		}
		// Check if the role's permissions contains the sought after permission
		if role.Permissions&perm != 0 {
			return true, nil
		}
	}
	return false, nil
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

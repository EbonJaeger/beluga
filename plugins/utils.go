package plugins

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

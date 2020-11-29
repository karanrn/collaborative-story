package helper

// Contains checks if value exists in the list
func Contains(key string, list []string) bool {
	for _, ut := range list {
		if key == ut {
			return true
		}
	}
	return false
}

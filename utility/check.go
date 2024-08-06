package utility

func InStringSlice(value string, strSlice []string) bool {
	for _, s := range strSlice {
		if value == s {
			return true
		}
	}
	return false
}
func InIntSlice(value int, strSlice []int) bool {
	for _, s := range strSlice {
		if value == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] == s {
			// Remove the string by slicing the slice
			slice = append(slice[:i], slice[i+1:]...)
			i-- // Decrement the index to re-check the current position
		}
	}
	return slice
}

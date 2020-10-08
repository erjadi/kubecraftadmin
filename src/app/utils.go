package main

// Contains checks for the occurence of a string in an array of strings
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Remove certain value from an array of strings
func Remove(a []string, x string) []string {
	if Contains(a, x) {
		for i, n := range a {
			if x == n {
				return append(a[:i], a[i+1:]...)
			}
		}
	}
	return a
}

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

package tools

import "sort"

// Returns true if the 'needle' exists in the 'haystack'
func SContains(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

// Are the elements of these slices identical?
// ORDER OF ITEMS CAN BE DIFFERENT (eg. [X, Y, Z] = [Z, X, Y])
// Params:
// - a, b: slices to compare
// Returns:
// - bool
func SameElements(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if len(a) == 0 {
		return true
	}

	// Copy slices so we don't overwrite the originals
	var copyA = make([]string, len(a))
	var copyB = make([]string, len(b))
	copy(copyA, a)
	copy(copyB, b)

	sort.Strings(copyA)
	sort.Strings(copyB)

	for i, v := range copyA {
		if copyB[i] != v {
			return false
		}
	}

	return true
}

// Only keep unique values in the string slice
// Params:
// - a: slice to filter
// Returns:
// - new slice with unique values
func UniqueStringSlice(a []string) []string {
	if len(a) == 0 {
		return a
	}

	uniqueItems := make(map[string]bool)

	for _, v := range a {
		uniqueItems[v] = true
	}

	filtered := make([]string, len(uniqueItems))
	i := 0
	for k := range uniqueItems {
		filtered[i] = k
		i++
	}

	return filtered
}

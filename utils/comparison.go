package utils

import "slices"

func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aItem := range a {
		if !slices.Contains(b, aItem) {
			return false
		}
	}

	return true
}

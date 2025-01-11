package engine

import (
	"github.com/chehsunliu/poker"
)

func CompareCardSlices(slice1, slice2 []poker.Card) bool {
	// Check if lengths are the same
	if len(slice1) != len(slice2) {
		return false
	}
	// Check each element for equality
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
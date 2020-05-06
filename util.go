package main

import "strings"

func floatOf(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func getBool(s string) bool {
	var bools = [...]string{"true", "t", "yes", "y", "false", "f", "no", "n"}
	// always compare lowercase
	s = strings.ToLower(s)
	for i, test := range bools {
		if test == s {
			// 0 >= i >= 3 is truthy, otherwise falsy
			return i < 4
		}
	}
	return false
}

func joinInt64(lo, hi uint32) int64 {
	return (int64(hi) << 32) + int64(lo)
}

package main

import (
	"strings"
)

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

func joinInt64(lo int32, hi uint32) int64 {
	// For some reason *Lo values might be negative on the serialized JSON received from NZBGet causing an error:
	// `json: cannot unmarshal number -1 into Go struct field temp.TotalSizeLo of type uint32`
	// See: https://forum.nzbget.net/viewtopic.php?t=3711
	// For this reason, we unmarshal them as signed then use the unsigned value
	ulo := uint32(lo)

	return (int64(hi) << 32) + int64(ulo)
}

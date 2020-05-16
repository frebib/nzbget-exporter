package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestServerVolumeUnmarshal(t *testing.T) {

	testCase := []byte(`
[
	{
		"ServerID": 0,
		"DataTime": 1589637733,
		"FirstDay": 17894,
		"TotalSizeLo": 330321954,
		"TotalSizeHi": 762,
		"TotalSizeMB": 3121467,
		"CustomSizeLo": 330321954,
		"CustomSizeHi": 762,
		"CustomSizeMB": 3121467,
		"CustomTime": 1546097163,
		"SecSlot": 13,
		"MinSlot": 2,
		"HourSlot": 16,
		"DaySlot": 504,
		"BytesPerSeconds": [
		{
			"SizeLo": 0,
			"SizeHi": 0,
			"SizeMB": 0
		}
		]
	},
	{
		"ServerID": 2,
		"DataTime": 1589637733,
		"FirstDay": 17894,
		"TotalSizeLo": -1059258586,
		"TotalSizeHi": 3,
		"TotalSizeMB": 15373,
		"CustomSizeLo": -1059258586,
		"CustomSizeHi": 3,
		"CustomSizeMB": 15373,
		"CustomTime": 1546097163,
		"SecSlot": 13,
		"MinSlot": 2,
		"HourSlot": 16,
		"DaySlot": 504,
		"BytesPerSeconds": [
		{
			"SizeLo": 0,
			"SizeHi": 0,
			"SizeMB": 0
		},
		{
			"SizeLo": 134212,
			"SizeHi": 20,
			"SizeMB": 81920
		}
		]
	}
]
`)

	out := []ServerVolume{}
	err := json.Unmarshal(testCase, &out)

	if err != nil {
		t.Errorf("Failed to unmarshal. Error: %s", err)
	}

	if !reflect.DeepEqual(out[0], ServerVolume{
		ID:         0,
		TotalBytes: 3273095401506, // 3273095401506/1024/1024 == 3121467 MB (TotalSizeMB)
	}) {
		t.Errorf("Unexpected ServerVolume[0]: %#v", out[0])
	}

	if !reflect.DeepEqual(out[1], ServerVolume{
		ID:         2,
		TotalBytes: 16120610598, // 16120610598/1024/1024 == 15373 MB (TotalSizeMB)
	}) {
		t.Errorf("Unexpected ServerVolume[1]: %#v", out[1])
	}
}

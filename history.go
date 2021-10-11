package main

import (
	"encoding/json"
	"time"
)

// https://nzbget.net/api/history

type (
	History struct {
		NZBID              uint64
		Name               string
		RemainingFileCount uint64
		RetryData          bool
		HistoryTime        time.Time `json:"-"`
		Status             string
		Log                []interface{}
		NZBName            string
		Kind               HistoryKind
		URL                string
		NZBFilename        string
		DestDir            string
		FinalDir           string
		Category           string
		ParStatus          ParStatus
		UnpackStatus       UnpackStatus
		MoveStatus         MoveStatus
		//ScriptStatus       ScriptStatus
		DeleteStatus     DeleteStatus
		MarkStatus       MarkStatus
		URLStatus        URLStatus
		FileSize         int64
		FileCount        uint64
		MinPostTime      time.Time `json:"-"`
		MaxPostTime      time.Time `json:"-"`
		TotalArticles    uint64
		SuccessArticles  uint64
		FailedArticles   uint64
		Health           uint64
		CriticalHealth   uint64
		DupeKey          string
		DupeScore        uint64
		DupeMode         string
		Deleted          bool
		DownloadedSize   int64
		DownloadTimeSec  uint64
		PostTotalTimeSec uint64
		ParTimeSec       uint64
		RepairTimeSec    uint64
		UnpackTimeSec    uint64
		MessageCount     uint64
		ExtraParBlocks   uint64
		Parameters       []Parameters
		ScriptStatuses   []ScriptStatus
		ServerStats      []ServerStats
	}

	Parameters struct {
		Name  string
		Value string
	}

	ScriptStatus struct {
		Name   string
		Status string
	}

	ServerStats struct {
		ServerID        int
		SuccessArticles int
		FailedArticles  int
	}
)

//go:generate enumerx -type=HistoryKind -trimprefix=Kind -json
type HistoryKind int

const (
	KindNZB HistoryKind = iota
	KindURL
	KindDUP
)

//go:generate enumerx -type=ParStatus -trimprefix=@type -transform=snake_upper -json
type ParStatus int

const (
	ParStatusNone ParStatus = iota
	ParStatusFailure
	ParStatusRepairPossible
	ParStatusSuccess
	ParStatusManual
)

//go:generate enumerx -type=UnpackStatus -trimprefix=@type -transform=snake_upper -json
type UnpackStatus int

const (
	UnpackStatusNone UnpackStatus = iota
	UnpackStatusFailure
	UnpackStatusSpace
	UnpackStatusPassword
	UnpackStatusSuccess
)

//go:generate enumerx -type=URLStatus -trimprefix=@type -transform=snake_upper -json
type URLStatus int

const (
	URLStatusNone URLStatus = iota
	URLStatusSuccess
	URLStatusFailure
	URLStatusScanSkipped
	URLStatusScanFailure
)

//go:generate enumerx -type=MoveStatus -trimprefix=@type -transform=snake_upper -json
type MoveStatus int

const (
	MoveStatusNone MoveStatus = iota
	MoveStatusSuccess
	MoveStatusFailure
)

//go:generate enumerx -type=DeleteStatus -trimprefix=@type -transform=snake_upper -json
type DeleteStatus int

const (
	DeleteStatusNone DeleteStatus = iota
	DeleteStatusManual
	DeleteStatusHealth
	DeleteStatusDupe
	DeleteStatusBad
	DeleteStatusScan
	DeleteStatusCopy
)

//go:generate enumerx -type=MarkStatus -trimprefix=@type -transform=snake_upper -json
type MarkStatus int

const (
	MarkStatusNone MarkStatus = iota
	MarkStatusGood
	MarkStatusBad
)

func (h *History) UnmarshalJSON(b []byte) error {
	// Unmarshal the struct as normal
	type resultClone History
	var clone = (*resultClone)(h)
	err := json.Unmarshal(b, clone)
	if err != nil {
		return err
	}

	type temp struct {
		FileSizeHi       uint32 `json:"FileSizeHi"`
		FileSizeLo       uint32 `json:"FileSizeLo"`
		DownloadedSizeHi uint32 `json:"DownloadedSizeHi"`
		DownloadedSizeLo uint32 `json:"DownloadedSizeLo"`

		HistoryTime int64 `json:"HistoryTime"`
		MinPostTime int64 `json:"MinPostTime"`
		MaxPostTime int64 `json:"MaxPostTime"`
	}

	values := temp{}
	err = json.Unmarshal(b, &values)
	if err != nil {
		return err
	}

	h.FileSize = joinInt64(values.FileSizeLo, values.FileSizeHi)
	h.DownloadedSize = joinInt64(values.DownloadedSizeLo, values.DownloadedSizeHi)

	h.HistoryTime = time.Unix(values.HistoryTime, 0)
	h.MinPostTime = time.Unix(values.MaxPostTime, 0)
	h.MaxPostTime = time.Unix(values.MaxPostTime, 0)

	return nil
}

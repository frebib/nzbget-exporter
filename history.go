package main

import (
	"encoding/json"
	"time"
)

// https://nzbget.net/api/history

type (
	History struct {
		NZBID              int
		Name               string
		RemainingFileCount int
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
		ParStatus          string
		ExParStatus        string
		UnpackStatus       string
		MoveStatus         string
		ScriptStatus       string
		DeleteStatus       string
		MarkStatus         string
		URLStatus          string
		FileSize           int64
		FileCount          int
		MinPostTime        time.Time `json:"-"`
		MaxPostTime        time.Time `json:"-"`
		TotalArticles      int
		SuccessArticles    int
		FailedArticles     int
		Health             int
		CriticalHealth     int
		DupeKey            string
		DupeScore          int
		DupeMode           string
		Deleted            bool
		DownloadedSize     int64
		DownloadTimeSec    int
		PostTotalTimeSec   int
		ParTimeSec         int
		RepairTimeSec      int
		UnpackTimeSec      int
		MessageCount       int
		ExtraParBlocks     int
		Parameters         []Parameters
		ScriptStatuses     []ScriptStatus
		ServerStats        []ServerStats
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

type HistoryKind int

//go:generate enumer -type=HistoryKind -trimprefix=Kind -json
const (
	KindNZB HistoryKind = iota
	KindURL
	KindDUP
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

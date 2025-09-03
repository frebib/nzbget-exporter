package main

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Version string          `json:"version"`
	Result  json.RawMessage `json:"result"`
}

type Status struct {
	ArticleCache   int64 `json:"-"`
	DaySize        int64 `json:"-"`
	DownloadedSize int64 `json:"-"`
	ForcedSize     int64 `json:"-"`
	FreeDiskSpace  int64 `json:"-"`
	MonthSize      int64 `json:"-"`
	RemainingSize  int64 `json:"-"`

	AverageDownloadRate int64 `json:"AverageDownloadRate"`
	DownloadLimit       int64 `json:"DownloadLimit"`
	DownloadRate        int64 `json:"DownloadRate"`
	DownloadTimeSec     int64 `json:"DownloadTimeSec"`
	ParJobCount         int64 `json:"ParJobCount"`
	PostJobCount        int64 `json:"PostJobCount"`
	ThreadCount         int64 `json:"ThreadCount"`
	URLCount            int64 `json:"UrlCount"`

	ServerPaused    bool `json:"ServerPaused"`
	DownloadPaused  bool `json:"DownloadPaused"`
	Download2Paused bool `json:"Download2Paused"`
	ServerStandBy   bool `json:"ServerStandBy"`
	PostPaused      bool `json:"PostPaused"`
	ScanPaused      bool `json:"ScanPaused"`
	QuotaReached    bool `json:"QuotaReached"`
	FeedActive      bool `json:"FeedActive"`

	ServerTime       time.Time     `json:"-"`
	ResumeTime       time.Time     `json:"-"`
	StartTime        time.Time     `json:"-"`
	QueueScriptCount int           `json:"QueueScriptCount"`
	NewsServers      []NewsServers `json:"NewsServers"`
}

type NewsServers struct {
	ID     int  `json:"ID"`
	Active bool `json:"Active"`
}

func (s *Status) UnmarshalJSON(b []byte) error {
	// Unmarshal the struct as normal
	type resultClone Status
	var clone *resultClone = (*resultClone)(s)
	err := json.Unmarshal(b, clone)
	if err != nil {
		return err
	}

	type temp struct {
		ArticleCacheHi   uint32 `json:"ArticleCacheHi"`
		ArticleCacheLo   uint32 `json:"ArticleCacheLo"`
		DaySizeHi        uint32 `json:"DaySizeHi"`
		DaySizeLo        uint32 `json:"DaySizeLo"`
		DownloadedSizeHi uint32 `json:"DownloadedSizeHi"`
		DownloadedSizeLo uint32 `json:"DownloadedSizeLo"`
		ForcedSizeHi     uint32 `json:"ForcedSizeHi"`
		ForcedSizeLo     uint32 `json:"ForcedSizeLo"`
		FreeDiskSpaceHi  uint32 `json:"FreeDiskSpaceHi"`
		FreeDiskSpaceLo  uint32 `json:"FreeDiskSpaceLo"`
		MonthSizeHi      uint32 `json:"MonthSizeHi"`
		MonthSizeLo      uint32 `json:"MonthSizeLo"`
		RemainingSizeHi  uint32 `json:"RemainingSizeHi"`
		RemainingSizeLo  uint32 `json:"RemainingSizeLo"`

		ServerTime int64 `json:"ServerTime"`
		ResumeTime int64 `json:"ResumeTime"`
		UptimeSec  int64 `json:"UpTimeSec"`
	}

	values := temp{}
	err = json.Unmarshal(b, &values)
	if err != nil {
		return err
	}

	s.ArticleCache = joinInt64(values.ArticleCacheLo, values.ArticleCacheHi)
	s.DaySize = joinInt64(values.DaySizeLo, values.DaySizeHi)
	s.DownloadedSize = joinInt64(values.DownloadedSizeLo, values.DownloadedSizeHi)
	s.ForcedSize = joinInt64(values.ForcedSizeLo, values.ForcedSizeHi)
	s.FreeDiskSpace = joinInt64(values.FreeDiskSpaceLo, values.FreeDiskSpaceHi)
	s.MonthSize = joinInt64(values.MonthSizeLo, values.MonthSizeHi)
	s.RemainingSize = joinInt64(values.RemainingSizeLo, values.RemainingSizeHi)

	s.ServerTime = time.Unix(values.ServerTime, 0)
	s.ResumeTime = time.Unix(values.ResumeTime, 0)
	s.StartTime = time.Now().Add(-(time.Second * time.Duration(values.UptimeSec)))

	return nil
}

type NZBGetConfig struct {
	AppBin            string
	AppDir            string
	AppendCategoryDir bool
	ArticleCache      int64
	ArticleInterval   int64
	ArticleRetries    int64
	ArticleTimeout    int64
	AuthorizedIP      string
	CertCheck         bool
	CertStore         string
	ConfigFile        string
	ConfigTemplate    string
	ContinuePartial   bool
	ControlIp         string
	ControlPort       string
	CrashDump         bool
	CrashTrace        bool
	CrcCheck          bool
	CursesGroup       bool
	CursesNzbName     bool
	CursesTime        bool
	DailyQuota        int
	DebugTarget       string
	DestDir           string
	DetailTarget      string
	DirectRename      bool
	DirectUnpack      bool
	DirectWrite       bool
	DiskSpace         int
	DownloadRate      int64
	DupeCheck         bool
	ErrorTarget       string
	EventInterval     string
	ExtCleanupDisk    string
	Extensions        string
	FeedHistory       int
	FileNaming        string
	FlushQueue        bool
	FormAuth          bool
	HealthCheck       string
	InfoTarget        string
	InterDir          string
	KeepHistory       int
	LockFile          string
	LogBuffer         int
	LogFile           string
	MainDir           string
	MonthlyQuota      int
	NzbCleanupDisk    bool
	NzbDir            string
	NzbDirFileAge     int
	NzbDirInterval    int
	NzbLog            bool
	OutputMode        string
	ParBuffer         int
	ParCheck          string
	ParIgnoreExt      []string
	ParPauseQueue     bool
	ParQuick          bool
	ParRename         bool
	ParRepair         bool
	ParScan           string
	ParThreads        int
	ParTimeLimit      string
	PostStrategy      string
	PropagationDelay  int
	QueueDir          string
	QuotaStartDay     int
	RarRename         bool
	RawArticle        bool
	RemoteTimeout     int
	ReorderFiles      bool
	RequiredDir       string
	RotateLog         int
	ScriptDir         string
	ScriptOrder       string
	ScriptPauseQueue  bool
	SecureCert        string
	SecureControl     bool
	SecureKey         string
	SecurePort        string
	SevenZipCmd       string
	ShellOverride     string
	SkipWrite         bool
	TempDir           string
	TimeCorrection    int
	UMask             string
	Unpack            bool
	UnpackCleanupDisk bool
	UnpackIgnoreExt   bool
	UnpackPassFile    string
	UnpackPauseQueue  bool
	UnrarCmd          string
	UpdateCheck       string
	UpdateInterval    int
	UrlConnections    int
	UrlForce          bool
	UrlInterval       int
	UrlRetries        int
	UrlTimeout        int
	Version           string
	WarningTarget     string
	WebDir            string
	WriteBuffer       int
	WriteLog          bool
	Server            []ConfigServer
}

type ConfigServer struct {
	Active      bool
	Cipher      string
	Connections int
	Encryption  bool
	Group       int
	Host        string
	IpVersion   string
	JoinGroup   bool
	Level       int
	Name        string
	Notes       string
	Optional    bool
	Port        uint16
	Retention   int
}

type ConfigCategory struct {
	Aliases    string
	DestDir    string
	Extensions string
	Name       string
	Unpack     string
}

func (c *NZBGetConfig) UnmarshalJSON(b []byte) error {
	var values []struct{ Name, Value string }
	err := json.Unmarshal(b, &values)
	if err != nil {
		return err
	}

	for _, val := range values {
		of := reflect.ValueOf(c)
		field := reflect.Indirect(of).FieldByName(val.Name)

		parts := strings.SplitN(val.Name, ".", 2)
		if len(parts) > 1 {
			regex := regexp.MustCompile(`^([^\d]+)(\d+)$`)
			groups := regex.FindStringSubmatch(parts[0])
			if len(groups) < 1 {
				continue
			}
			structName := groups[1]
			fieldName := parts[1]
			key, _ := strconv.Atoi(groups[2]) // Second group is the number
			key--                             // zero-indexed, so 1 becomes 0

			slice := reflect.Indirect(of).FieldByName(structName)
			if slice.Kind() != reflect.Slice {
				continue
			}
			if key >= slice.Cap() {
				nCap := slice.Cap() * 2
				if nCap == 0 {
					nCap = 2
				}
				for nCap < key {
					nCap *= 2
				}
				// make a new slice, and copy existing values in
				// golang doesn't expose the "grow" functionality, so we
				// do it ourselves to be able to access a particular index
				bigger := reflect.MakeSlice(slice.Type(), nCap, nCap)
				reflect.Copy(bigger, slice)
				slice.Set(bigger)
			}
			nthElem := slice.Index(key)
			// TODO: Handle nthElem being nil-able type
			field = reflect.Indirect(nthElem).FieldByName(fieldName)
		}

		reflectInto(field, val.Value)
	}

	return nil
}

func reflectInto(v reflect.Value, str string) {
	if !v.CanSet() {
		return
	}
	switch v.Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		uintVal, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return
		}
		v.SetUint(uintVal)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		intVal, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return
		}
		v.SetInt(intVal)
	case reflect.Bool:
		// if string is not bool-y then default false
		v.SetBool(getBool(str))
	case reflect.String:
		v.SetString(str)
	case reflect.Slice:
		strs := strings.Split(str, ",")
		v.Set(reflect.MakeSlice(v.Type(), len(strs), len(strs)))
		for i := 0; i < v.Len(); i++ {
			index := v.Index(i)
			reflectInto(index, strings.TrimSpace(strs[i]))
		}
	default:
		return
	}
}

type ServerVolume struct {
	ID                  int   `json:"-"`
	TotalBytes          int64 `json:"-"`
	TotalArticleSuccess int   `json:"-"`
	TotalArticleFailed int   `json:"-"`
}

func (v *ServerVolume) UnmarshalJSON(b []byte) error {

	type ArticlePerDay struct {
		Success int `json:"Success"`
		Failed  int `json:"Failed"`
	}

	type temp struct {
		ServerID        int             `json:"ServerID"`
		TotalSizeLo     uint32          `json:"TotalSizeLo"`
		TotalSizeHi     uint32          `json:"TotalSizeHi"`
		ArticlesPerDays []ArticlePerDay `json:"ArticlesPerDays"`
	}

	values := temp{}
	err := json.Unmarshal(b, &values)
	if err != nil {
		return err
	}

	v.ID = values.ServerID
	v.TotalBytes = joinInt64(values.TotalSizeLo, values.TotalSizeHi)

	var totalSuccess, totalFailed int
	for _, day := range values.ArticlesPerDays {
		totalSuccess += day.Success
		totalFailed += day.Failed
	}

	v.TotalArticleSuccess = totalSuccess
	v.TotalArticleFailed = totalFailed

	return nil
}

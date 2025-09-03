package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	prom "github.com/prometheus/client_golang/prometheus"
)

type NZBGetCollector struct {
	Config *ExporterConfig

	version *prom.Desc

	articleCache    *prom.Desc
	diskSpaceFree   *prom.Desc
	diskSpaceMin    *prom.Desc
	downloadLimit   *prom.Desc
	downloadPaused  *prom.Desc
	downloadTimeSec *prom.Desc
	downloadedSize  *prom.Desc
	forcedSize      *prom.Desc
	postJobCount    *prom.Desc
	postPaused      *prom.Desc
	quotaDay        *prom.Desc
	quotaMonth      *prom.Desc
	quotaReached    *prom.Desc
	remainingSize   *prom.Desc
	resumeTime      *prom.Desc
	scanPaused      *prom.Desc
	serverStandBy   *prom.Desc
	startTime       *prom.Desc
	threadCount     *prom.Desc
	urlCount        *prom.Desc

	newsServerActive         *prom.Desc
	newsServerBytes          *prom.Desc
	newsServerArticleSuccess *prom.Desc
	newsServerArticleFailed *prom.Desc

	historyCategoryCount       *prom.Desc
	historyFileSizeBytes       *prom.Desc
	historyFileCount           *prom.Desc
	historyRemainingFileCount  *prom.Desc
	historyArticleCount        *prom.Desc
	historySuccessArticleCount *prom.Desc
	historyFailedArticleCount  *prom.Desc
	historyDownloadTime        *prom.Desc
	historyDownloadSizeBytes   *prom.Desc
	historyPostTime            *prom.Desc
	historyParTime             *prom.Desc
	historyRepairTime          *prom.Desc
	historyUnpackTime          *prom.Desc
	historyStatusCount         *prom.Desc
	historyParStatusCount      *prom.Desc
	historyUnpackStatusCount   *prom.Desc
}

func NewNZBGetCollector(config *ExporterConfig) *NZBGetCollector {
	ns := config.Namespace

	return &NZBGetCollector{
		Config: config,

		version: prom.NewDesc(
			prom.BuildFQName(ns, "", "version"),
			"always 1. label 'version' contains nzbget server version",
			[]string{"version"}, nil,
		),

		articleCache: prom.NewDesc(
			prom.BuildFQName(ns, "article_cache", "bytes"),
			"Current usage of article cache",
			nil, nil,
		),
		diskSpaceFree: prom.NewDesc(
			prom.BuildFQName(ns, "disk", "free_bytes"),
			"Free disk space on 'DestDir'",
			nil, nil,
		),
		diskSpaceMin: prom.NewDesc(
			prom.BuildFQName(ns, "disk", "min_bytes"),
			"Disk space limit before pausing the download queue",
			nil, nil,
		),
		downloadLimit: prom.NewDesc(
			prom.BuildFQName(ns, "download", "limit"),
			"Current download limit, in bytes per second",
			nil, nil,
		),
		downloadPaused: prom.NewDesc(
			prom.BuildFQName(ns, "download", "paused"),
			"1 if the download queue is paused, 0 otherwise",
			nil, nil,
		),
		downloadTimeSec: prom.NewDesc(
			prom.BuildFQName(ns, "download", "time_seconds"),
			"Server download time in seconds",
			nil, nil,
		),
		downloadedSize: prom.NewDesc(
			prom.BuildFQName(ns, "downloaded", "total_bytes"),
			"Total data downloaded since server start",
			nil, nil,
		),
		forcedSize: prom.NewDesc(
			prom.BuildFQName(ns, "forced", "bytes"),
			"Remaining size of entries with FORCE priority",
			nil, nil,
		),
		postJobCount: prom.NewDesc(
			prom.BuildFQName(ns, "post", "job_count"),
			"Number of Par-Jobs or Post-processing script jobs in the post-processing queue",
			nil, nil,
		),
		postPaused: prom.NewDesc(
			prom.BuildFQName(ns, "post", "active"),
			"1 if post-processor queue is currently active, 0 if paused",
			nil, nil,
		),
		quotaDay: prom.NewDesc(
			prom.BuildFQName(ns, "quota", "day_bytes"),
			"Daily quota in bytes", nil, nil,
		),
		quotaMonth: prom.NewDesc(
			prom.BuildFQName(ns, "quota", "month_bytes"),
			"Monthly quota in bytes", nil, nil,
		),
		quotaReached: prom.NewDesc(
			prom.BuildFQName(ns, "quota", "reached"),
			"1 if quota has been hit, 0 otherwise", nil, nil,
		),
		remainingSize: prom.NewDesc(
			prom.BuildFQName(ns, "queue", "remaining_bytes"),
			"Remaining size of all entries in download queue",
			nil, nil,
		),
		resumeTime: prom.NewDesc(
			prom.BuildFQName(ns, "resume", "time"),
			"Time to resume if set with method \"scheduleresume\"",
			nil, nil,
		),
		scanPaused: prom.NewDesc(
			prom.BuildFQName(ns, "scan", "active"),
			"1 if the scanning of incoming nzb-directory is currently active, 0 if paused",
			nil, nil,
		),
		serverStandBy: prom.NewDesc(
			prom.BuildFQName(ns, "", "standby"),
			"1 if no downloads in progress (server paused or all jobs completed), otherwise 0 if there are currently downloads running",
			nil, nil,
		),
		startTime: prom.NewDesc(
			prom.BuildFQName(ns, "start_time", "seconds"),
			"Server start time, in unixtime",
			nil, nil,
		),
		threadCount: prom.NewDesc(
			prom.BuildFQName(ns, "thread", "count"),
			"Number of threads running",
			nil, nil,
		),
		urlCount: prom.NewDesc(
			prom.BuildFQName(ns, "url", "count"),
			"Number of URLs in the URL-queue (including current file)",
			nil, nil,
		),

		newsServerActive: prom.NewDesc(
			prom.BuildFQName(ns, "news_server", "active"),
			"News server used for obtaining articles, 1 if active",
			[]string{"id", "server"}, nil,
		),
		newsServerBytes: prom.NewDesc(
			prom.BuildFQName(ns, "news_server", "total_bytes"),
			"Total bytes downloaded from this news server",
			[]string{"id", "server"}, nil,
		),
		newsServerArticleSuccess: prom.NewDesc(
			prom.BuildFQName(ns, "news_server", "total_article_success"),
			"Total successful articles from this news server",
			[]string{"id", "server"}, nil,
		),
		newsServerArticleFailed: prom.NewDesc(
			prom.BuildFQName(ns, "news_server", "total_article_failed"),
			"Total failed articles from this news server",
			[]string{"id", "server"}, nil,
		),
		historyCategoryCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_category", "count"),
			"Number of history items in each category",
			[]string{"category"}, nil,
		),
		historyFileSizeBytes: prom.NewDesc(
			prom.BuildFQName(ns, "history_file_size", "total_bytes"),
			"Total bytes of all files in history",
			nil, nil,
		),
		historyFileCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_file", "count"),
			"Number of files in history",
			nil, nil,
		),
		historyRemainingFileCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_file", "remaining_count"),
			"Number of remaining files parked in history",
			nil, nil,
		),
		historyArticleCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_article", "count"),
			"Number of articles in history",
			nil, nil,
		),
		historySuccessArticleCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_article", "success_count"),
			"Number of successful articles in history",
			nil, nil,
		),
		historyFailedArticleCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_article", "failed_count"),
			"Number of failed articles in history",
			nil, nil,
		),
		historyDownloadTime: prom.NewDesc(
			prom.BuildFQName(ns, "history_download", "time_seconds"),
			"Download time in seconds",
			nil, nil,
		),
		historyDownloadSizeBytes: prom.NewDesc(
			prom.BuildFQName(ns, "history_download", "size_bytes"),
			"Total downloaded size in history, in bytes",
			nil, nil,
		),
		historyPostTime: prom.NewDesc(
			prom.BuildFQName(ns, "history_post_process", "time_seconds"),
			"Total post-processing time in seconds in history",
			nil, nil,
		),
		historyParTime: prom.NewDesc(
			prom.BuildFQName(ns, "history_par", "time_seconds"),
			"Total par-check time in seconds in history",
			nil, nil,
		),
		historyRepairTime: prom.NewDesc(
			prom.BuildFQName(ns, "history_par_repair", "time_seconds"),
			"Par-repair time in seconds in history",
			nil, nil,
		),
		historyUnpackTime: prom.NewDesc(
			prom.BuildFQName(ns, "history_unpack", "time_seconds"),
			"Unpack time in seconds in history",
			nil, nil,
		),
		historyStatusCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_status", "count"),
			"Number of history items per status",
			[]string{"reason", "status"}, nil,
		),
		historyParStatusCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_par_status", "count"),
			"Number of history items per par status",
			[]string{"status"}, nil,
		),
		historyUnpackStatusCount: prom.NewDesc(
			prom.BuildFQName(ns, "history_unpack_status", "count"),
			"Number of history items per unpack status",
			[]string{"status"}, nil,
		),
	}
}

func (c *NZBGetCollector) Collect(metrics chan<- prom.Metric) {
	var config NZBGetConfig
	var status Status
	var version string
	var volume []ServerVolume
	var history []History

	var wg sync.WaitGroup
	wg.Add(4)

	// Wait for config separately as multiple gothreads require it
	var cfgWg sync.WaitGroup
	var cfgErr = false
	cfgWg.Add(1)

	go func() {
		defer cfgWg.Done()
		err := c.getApi("config", &config)
		if err != nil {
			log.WithError(err).Error("api get config")
			metrics <- prom.NewInvalidMetric(prom.NewInvalidDesc(err), err)
			cfgErr = true
			return
		}
		metrics <- prom.MustNewConstMetric(c.diskSpaceMin, prom.GaugeValue, float64(config.DiskSpace*1024*1024))
	}()

	go func() {
		defer wg.Done()

		err := c.getApi("version", &version)
		if err != nil {
			log.WithError(err).Error("api get version")
			metrics <- prom.NewInvalidMetric(prom.NewInvalidDesc(err), err)
			return
		}
		metrics <- prom.MustNewConstMetric(c.version, prom.GaugeValue, 1, version)
	}()

	go func() {
		defer wg.Done()

		err := c.getApi("status", &status)
		if err != nil {
			log.WithError(err).Error("api get status")
			metrics <- prom.NewInvalidMetric(prom.NewInvalidDesc(err), err)
			return
		}
		metrics <- prom.MustNewConstMetric(c.articleCache, prom.GaugeValue, float64(status.ArticleCache))
		metrics <- prom.MustNewConstMetric(c.diskSpaceFree, prom.GaugeValue, float64(status.FreeDiskSpace))
		metrics <- prom.MustNewConstMetric(c.downloadLimit, prom.GaugeValue, float64(status.DownloadLimit))
		metrics <- prom.MustNewConstMetric(c.downloadPaused, prom.GaugeValue, floatOf(status.DownloadPaused))
		metrics <- prom.MustNewConstMetric(c.downloadTimeSec, prom.GaugeValue, float64(status.DownloadTimeSec))
		metrics <- prom.MustNewConstMetric(c.downloadedSize, prom.CounterValue, float64(status.DownloadedSize))
		metrics <- prom.MustNewConstMetric(c.forcedSize, prom.GaugeValue, float64(status.ForcedSize))
		metrics <- prom.MustNewConstMetric(c.postJobCount, prom.GaugeValue, float64(status.PostJobCount))
		metrics <- prom.MustNewConstMetric(c.postPaused, prom.GaugeValue, floatOf(status.PostPaused))
		metrics <- prom.MustNewConstMetric(c.quotaDay, prom.GaugeValue, float64(status.DaySize))
		metrics <- prom.MustNewConstMetric(c.quotaMonth, prom.GaugeValue, float64(status.MonthSize))
		metrics <- prom.MustNewConstMetric(c.quotaReached, prom.GaugeValue, floatOf(status.QuotaReached))
		metrics <- prom.MustNewConstMetric(c.remainingSize, prom.GaugeValue, float64(status.RemainingSize))
		metrics <- prom.MustNewConstMetric(c.resumeTime, prom.GaugeValue, float64(status.ResumeTime.Unix()))
		metrics <- prom.MustNewConstMetric(c.scanPaused, prom.GaugeValue, floatOf(status.ScanPaused))
		metrics <- prom.MustNewConstMetric(c.serverStandBy, prom.GaugeValue, floatOf(status.ServerStandBy))
		metrics <- prom.MustNewConstMetric(c.startTime, prom.GaugeValue, float64(status.StartTime.Unix()))
		metrics <- prom.MustNewConstMetric(c.threadCount, prom.GaugeValue, float64(status.ThreadCount))
		metrics <- prom.MustNewConstMetric(c.urlCount, prom.GaugeValue, float64(status.URLCount))

		cfgWg.Wait()
		if cfgErr {
			return
		}
		for _, srv := range status.NewsServers {
			idx := srv.ID
			id := fmt.Sprintf("%d", srv.ID)
			name := config.Server[idx-1].Name
			active := floatOf(srv.Active)

			metrics <- prom.MustNewConstMetric(c.newsServerActive, prom.GaugeValue, active, id, name)
		}
	}()

	go func() {
		defer wg.Done()

		err := c.getApi("servervolumes", &volume)
		if err != nil {
			log.WithError(err).Error("api get servervolumes")
			metrics <- prom.NewInvalidMetric(prom.NewInvalidDesc(err), err)
			return
		}

		cfgWg.Wait()
		if cfgErr {
			return
		}
		// https://nzbget.net/api/servervolumes
		// NOTE: The first record (serverid=0) are totals for all servers
		for _, srv := range volume {
			if srv.ID == 0 {
				continue
			}
			idx := srv.ID
			id := fmt.Sprintf("%d", srv.ID)
			name := config.Server[idx-1].Name
			bytes := float64(volume[idx].TotalBytes)

			metrics <- prom.MustNewConstMetric(c.newsServerBytes, prom.GaugeValue, bytes, id, name)

			metrics <- prom.MustNewConstMetric(c.newsServerArticleSuccess, prom.CounterValue, float64(volume[idx].TotalArticleSuccess), id, name)
			metrics <- prom.MustNewConstMetric(c.newsServerArticleFailed, prom.CounterValue, float64(volume[idx].TotalArticleFailed), id, name)

		}
	}()

	go func() {
		defer wg.Done()

		err := c.getApi("history", &history)
		if err != nil {
			log.WithError(err).Error("api get history")
			metrics <- prom.NewInvalidMetric(prom.NewInvalidDesc(err), err)
			return
		}

		var (
			fileSize int64
			fileCount,
			remainingCount,
			articleCount,
			articleSuccessCount,
			articleFailureCount,
			downloadTime uint64
			downloadSize int64
			postTime,
			parTime,
			repairTime,
			unpackTime uint64

			categories   = map[string]uint64{}
			statuses     = map[string]map[string]uint64{}
			parStatus    = map[string]uint64{}
			unpackStatus = map[string]uint64{}
		)

		for _, hi := range history {
			fileSize += hi.FileSize
			fileCount += hi.FileCount
			remainingCount += hi.RemainingFileCount
			articleCount += hi.TotalArticles
			articleSuccessCount += hi.SuccessArticles
			articleFailureCount += hi.FailedArticles
			downloadTime += hi.DownloadTimeSec
			downloadSize += hi.DownloadedSize
			postTime += hi.PostTotalTimeSec
			parTime += hi.ParTimeSec
			repairTime += hi.ParTimeSec
			unpackTime += hi.UnpackTimeSec

			categories[hi.Category]++
			parStatus[strings.ToLower(hi.ParStatus.String())]++
			unpackStatus[strings.ToLower(hi.UnpackStatus.String())]++

			// status is 'status/reason' such as 'success/health'
			parts := strings.Split(strings.ToLower(hi.Status), "/")
			status := parts[0]
			reason := parts[1]
			if statuses[status] == nil {
				statuses[status] = map[string]uint64{}
			}
			statuses[status][reason]++
		}

		metrics <- prom.MustNewConstMetric(c.historyFileSizeBytes, prom.CounterValue, float64(fileSize))
		metrics <- prom.MustNewConstMetric(c.historyFileCount, prom.CounterValue, float64(fileCount))
		metrics <- prom.MustNewConstMetric(c.historyRemainingFileCount, prom.CounterValue, float64(remainingCount))
		metrics <- prom.MustNewConstMetric(c.historyArticleCount, prom.CounterValue, float64(articleCount))
		metrics <- prom.MustNewConstMetric(c.historySuccessArticleCount, prom.CounterValue, float64(articleSuccessCount))
		metrics <- prom.MustNewConstMetric(c.historyFailedArticleCount, prom.CounterValue, float64(articleFailureCount))
		metrics <- prom.MustNewConstMetric(c.historyDownloadTime, prom.CounterValue, float64(downloadTime))
		metrics <- prom.MustNewConstMetric(c.historyDownloadSizeBytes, prom.CounterValue, float64(downloadSize))
		metrics <- prom.MustNewConstMetric(c.historyPostTime, prom.CounterValue, float64(postTime))
		metrics <- prom.MustNewConstMetric(c.historyParTime, prom.CounterValue, float64(parTime))
		metrics <- prom.MustNewConstMetric(c.historyRepairTime, prom.CounterValue, float64(repairTime))
		metrics <- prom.MustNewConstMetric(c.historyUnpackTime, prom.CounterValue, float64(unpackTime))
		sendConstMapMetric(metrics, c.historyCategoryCount, prom.CounterValue, categories)
		sendConstMapMapMetric(metrics, c.historyStatusCount, prom.CounterValue, statuses)
		sendConstMapMetric(metrics, c.historyParStatusCount, prom.CounterValue, parStatus)
		sendConstMapMetric(metrics, c.historyUnpackStatusCount, prom.CounterValue, unpackStatus)
	}()

	wg.Wait()
	cfgWg.Wait()
}

func (c *NZBGetCollector) getApi(endpoint string, out interface{}) error {
	// Remove right-trailing slashes, otherwise NZBGet will 404
	host := strings.TrimRight(c.Config.Host, "/")

	u, err := url.Parse(host + "/jsonrpc/" + endpoint)
	if err != nil {
		return err
	}
	log.WithField("url", u.String()).Debug("GET api")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	if c.Config.Username != "" && c.Config.Password != "" {
		req.SetBasicAuth(c.Config.Username, c.Config.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nzbget api response %d %s",
			resp.StatusCode, http.StatusText(resp.StatusCode),
		)
	}
	var response = new(Response)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	return json.Unmarshal(response.Result, out)
}

func sendConstMapMetric(metrics chan<- prom.Metric, desc *prom.Desc, valueType prom.ValueType, values map[string]uint64, labelValues ...string) {
	for key, value := range values {
		labels := append([]string{key}, labelValues...)
		metrics <- prom.MustNewConstMetric(desc, valueType, float64(value), labels...)
	}
}
func sendConstMapMapMetric(metrics chan<- prom.Metric, desc *prom.Desc, valueType prom.ValueType, values map[string]map[string]uint64, labelValues ...string) {
	for key, inner := range values {
		labels := append([]string{key}, labelValues...)
		sendConstMapMetric(metrics, desc, valueType, inner, labels...)
	}
}

func (c *NZBGetCollector) Describe(descr chan<- *prom.Desc) {
	descr <- c.articleCache
	descr <- c.diskSpaceFree
	descr <- c.diskSpaceMin
	descr <- c.downloadLimit
	descr <- c.downloadPaused
	descr <- c.downloadTimeSec
	descr <- c.downloadedSize
	descr <- c.forcedSize
	descr <- c.postJobCount
	descr <- c.postPaused
	descr <- c.quotaDay
	descr <- c.quotaMonth
	descr <- c.quotaReached
	descr <- c.remainingSize
	descr <- c.resumeTime
	descr <- c.scanPaused
	descr <- c.serverStandBy
	descr <- c.startTime
	descr <- c.threadCount
	descr <- c.urlCount

	descr <- c.newsServerActive
	descr <- c.newsServerBytes

	descr <- c.historyCategoryCount
	descr <- c.historyFileSizeBytes
	descr <- c.historyFileCount
	descr <- c.historyRemainingFileCount
	descr <- c.historyArticleCount
	descr <- c.historySuccessArticleCount
	descr <- c.historyFailedArticleCount
	descr <- c.historyDownloadTime
	descr <- c.historyDownloadSizeBytes
	descr <- c.historyPostTime
	descr <- c.historyParTime
	descr <- c.historyRepairTime
	descr <- c.historyUnpackTime
	descr <- c.historyStatusCount
	descr <- c.historyParStatusCount
	descr <- c.historyUnpackStatusCount
}

var _ prom.Collector = &NZBGetCollector{}

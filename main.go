package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	Version string = "unknown"
	log     *logrus.Entry
)

func init() {
	log = logrus.WithContext(context.Background())
	log.Logger.SetOutput(os.Stderr)
	log.Logger.Formatter = &prefixed.TextFormatter{
		FullTimestamp:  true,
		QuoteCharacter: "'",
	}
}

func main() {
	var config ExporterConfig
	parser := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)
	parser.Groups()[0].ShortDescription = "Options"
	_, err := parser.Parse()
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			fmt.Fprintf(os.Stderr, "NZBGet Exporter (version %s)\n\n", Version)
			parser.WriteHelp(os.Stderr)
			os.Exit(0)
		} else {
			log.WithError(err).
				Fatal("parse flags")
		}
	}

	log.Info("nzbget-exporter version " + Version)

	// Collect metrics for the provided backup provider
	collector := NewNZBGetCollector(&config)
	prom.MustRegister(collector)

	var version string
	err = collector.getApi("version", &version)
	if err != nil {
		log.WithError(err).Warn("failed to get nzbget version")
	} else {
		log.Infof("nzbget version %s", version)
	}

	promHandler := promhttp.Handler()
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		log.WithField("remote", r.RemoteAddr).
			Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		promHandler.ServeHTTP(w, r)
	}

	log.Info("serving metrics at " + config.Listen)

	http.Handle("/metrics", handler)
	err = http.ListenAndServe(config.Listen, nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Panic("listenandserve")
	}
}

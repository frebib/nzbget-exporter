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
	"github.com/x-cray/logrus-prefixed-formatter"
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
	log.Info("nzbget-exporter version " + Version)

	var config Config
	parser := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)
	parser.Group.LongDescription = "NZBGet Exporter Options"
	_, err := parser.Parse()
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
			os.Exit(0)
		} else {
			log.WithError(err).
				Fatal("parse flags")
		}
	}

	// Collect metrics for the provided backup provider
	collector := NewNZBGetCollector(config.Namespace)
	collector.Config = config.NZBGet
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

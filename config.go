package main

type ExporterConfig struct {
	LogLevel  string `long:"log-level" description:"log verbosity level (trace, debug, info, warn, error, fatal)" env:"LOG_LEVEL" default:"info"`
	Namespace string `long:"namespace" description:"metric name prefix" default:"nzbget" env:"NZBGET_METRIC_NAMESPACE"`
	Listen    string `short:"l" long:"listen" description:"host:port to listen on" default:":9452" env:"NZBGET_LISTEN"`
	Host      string `short:"h" long:"host" description:"nzbget host to export metrics for" required:"true" env:"NZBGET_HOST"`
	Username  string `short:"u" long:"username" description:"nzbget username for basicauth" env:"NZBGET_USERNAME"`
	Password  string `short:"p" long:"password" description:"nzbget password for basicauth" env:"NZBGET_PASSWORD"`
}

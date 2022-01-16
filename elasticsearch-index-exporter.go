package main

import (
	"flag"
	"fmt"
	"os"
	"prometheus-elasticsearch-index-exporter/config"
	"prometheus-elasticsearch-index-exporter/exporter"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	printVersion           = flag.Bool("version", false, "Print the prometheus-elasticsearch-index-exporter version")
	configPath             = flag.String("config", "", "Path to the config file. Try '-config ./example/config.yml' to get started.")
	showConfig             = flag.Bool("showconfig", true, "Print the current configuration to the console. Example: 'elasticsearch-index-exporter -showconfig -config ./example/config.yml'")
	disableExporterMetrics = flag.Bool("disable-exporter-metrics", false, "If this flag is set, the metrics about the exporter itself (go_*, process_*, promhttp_*) will be excluded from /metrics")
)

func main() {
	flag.Parse()
	if *printVersion {
		fmt.Println("version: 0.0.1-dev")
		return
	}

	validateCommandLineOrExit()

	cfg, warn, err := config.LoadConfigFile(*configPath)
	if len(warn) > 0 && !*showConfig {
		// warning is suppressed when '-showconfig' is used
		fmt.Fprintf(os.Stderr, "%v\n", warn)
	}
	exitOnError(err)
	if *showConfig {
		fmt.Printf("%+v\n", cfg)
		//return
	}

	// gather up the handlers with which to start the webserver
	var httpHandlers []exporter.HttpServerPathHandler
	httpHandlers = append(httpHandlers, exporter.HttpServerPathHandler{
		Path:    cfg.Server.Path,
		Handler: promhttp.Handler(),
	})
	// metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	// if !*disableExporterMetrics {
	// 	metricsHandler = promhttp.InstrumentMetricHandler(registry, metricsHandler)
	// }
	// httpHandlers = append(httpHandlers, exporter.HttpServerPathHandler{
	// 	Path:    cfg.Server.Path,
	// 	Handler: metricsHandler,
	// })

	fmt.Print(startMsg(cfg, httpHandlers))
	serverErrors := startServer(cfg.Server, httpHandlers)

	retentionTicker := time.NewTicker(cfg.Global.RetentionCheckInterval)

	for {
		select {
		case err := <-serverErrors:
			exitOnError(fmt.Errorf("server error: %v", err.Error()))
		case <-retentionTicker.C:
			fmt.Println("Retention Ticker")
			// for _, metric := range metrics {
			// 	err = metric.ProcessRetention()
			// 	if err != nil {
			// 		fmt.Fprintf(os.Stderr, "WARNING: error while processing retention on metric %v: %v", metric.Name(), err)
			// 		nErrorsByMetric.WithLabelValues(metric.Name()).Inc()
			// 	}
			// }
		}
	}
}

func startMsg(cfg *config.Config, httpHandlers []exporter.HttpServerPathHandler) string {
	host := "localhost"
	if len(cfg.Server.Host) > 0 {
		host = cfg.Server.Host
	} else {
		hostname, err := os.Hostname()
		if err == nil {
			host = hostname
		}
	}

	var sb strings.Builder
	baseUrl := fmt.Sprintf("%v://%v:%v", cfg.Server.Protocol, host, cfg.Server.Port)
	sb.WriteString("Starting server on")
	for _, httpHandler := range httpHandlers {
		sb.WriteString(fmt.Sprintf(" %v%v", baseUrl, httpHandler.Path))
	}
	sb.WriteString("\n")
	return sb.String()
}

func validateCommandLineOrExit() {
	if len(*configPath) == 0 {
		if *showConfig {
			fmt.Fprint(os.Stderr, "Usage: elasticsearch-index-exporter -showconfig -config <path>\n")
		} else {
			fmt.Fprint(os.Stderr, "Usage: elasticsearch-index-exporter -config <path>\n")
		}
		os.Exit(-1)
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(-1)
	}
}

func startServer(cfg config.ServerConfig, httpHandlers []exporter.HttpServerPathHandler) chan error {
	serverErrors := make(chan error)
	go func() {
		serverErrors <- exporter.RunHttpServer(cfg.Host, cfg.Port, httpHandlers)
	}()
	return serverErrors
}

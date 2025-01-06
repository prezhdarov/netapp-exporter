package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prezhdarov/prometheus-exporter/config"
	"github.com/prezhdarov/prometheus-exporter/exporter"

	ONTAPI "netapp-exporter/netapp/api"
	ontapCollectors "netapp-exporter/netapp/collectors"

	"github.com/prometheus/common/promslog"
	"github.com/prometheus/exporter-toolkit/web"
)

const (
	exporterName = "NetApp ONTAP REST API Exporter"
	namespace    = "netapp"
)

var (
	listenAddress          = flag.String("http.address", ":9169", "Address and port to listen for http connections")
	maxRequests            = flag.Int("prom.maxRequests", 20, "Maximum number of parallel scrape requests. Use 0 to disable.")
	disableExporterTarget  = flag.Bool("disable.exporter.target", false, "Disable default target for /metrics path.")
	disableExporterMetrics = flag.Bool("disable.exporter.metrics", true, "Disable exporter metrics in /metrics path. Always enabled if /metrics target disabled")

	logLevel  = flag.String("log.level", "debug", "Log Level minimums. Available options are: debug,info,warn and error")
	logFormat = flag.String("log.format", "logfmt", "Log output format. Available options are: logfmt and json")
)

func usage() {
	const s = `
netapp-exporter collects metrics data from NetApp ONTAP Cluster. 
`
	config.Usage(s)
}

func main() {

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = usage
	config.Parse()

	logger := promslog.New(config.SetLogger(logFormat, logLevel))

	logger.Debug("exporter target is", fmt.Sprintf("%v", *disableExporterTarget), nil)

	ONTAPI.Load(logger)

	ontapCollectors.Load(logger)

	http.Handle("/metrics", exporter.CreateHandler(!*disableExporterMetrics, *disableExporterTarget, *maxRequests, namespace, logger))
	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		exporter.CreateHandleFunc(w, r, namespace, "gateways", logger)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<head><title>NetApp ONTAP REST API Exporter</title></head>
			<body>
			<h1>NetApp ONTAP REST API Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			<p><a href="/probe">Probe</a></p>
			</body>
			</html>`))
	})

	logger.Info("msg", "listening on", "address", *listenAddress, nil)

	server := &http.Server{}

	if err := web.ListenAndServe(server, config.WebConfig(listenAddress), logger); err != nil {
		logger.Error("err", fmt.Sprintf("%s", err), nil)
		os.Exit(1)
	}

}

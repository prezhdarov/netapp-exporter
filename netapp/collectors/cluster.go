package ontapCollectors

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/go-kit/log"
	"github.com/prezhdarov/prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	clusterSubsystem = "cluster"
)

var clusterCollectorFlag = flag.Bool(fmt.Sprintf("collector.%s", clusterSubsystem), collector.DefaultEnabled, fmt.Sprintf("Enable the %s collector (default: %v)", clusterSubsystem, collector.DefaultEnabled))

type clusterCollector struct {
	logger log.Logger
}

func init() {
	collector.RegisterCollector("cluster", clusterCollectorFlag, NewClusterCollector)
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewClusterCollector(logger log.Logger) (collector.Collector, error) {
	return &clusterCollector{logger}, nil
}

func (c *clusterCollector) Update(ch chan<- prometheus.Metric, namespace string, clientAPI collector.ClientAPI, loginData map[string]interface{}, params map[string]string) error {

	var cluster Cluster

	extraConfig := make(map[string]interface{}, 0)

	extraConfig["api"] = "/api/cluster/?fields=metric,version"
	body, err := clientAPI.Get(loginData, extraConfig, c.logger)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body.(*[]byte), &cluster)
	if err != nil {
		return err
	}

	labels := map[string]string{"na_cluster": loginData["target"].(string)}

	labels["info"] = strings.Trim(strings.TrimPrefix(strings.Split(cluster.Full, ":")[0], "NetApp Release"), " ")

	ch <- prometheus.NewMetricWithTimestamp(
		cluster.Metric.ToTimestamp(c.logger), prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, aggregateSubsystem, "info"),
				"ONTAP Cluster Inormation",
				nil, labels,
			),
			prometheus.GaugeValue, 1.0,
		),
	)

	return nil
}

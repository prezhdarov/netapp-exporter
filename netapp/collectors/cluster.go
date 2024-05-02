package ontapCollectors

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/prezhdarov/prometheus-exporter/collector"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	clusterSubsystem = "cluster"
)

var clusterCollectorFlag = flag.Bool("cluster.collector", collector.DefaultEnabled, fmt.Sprintf("Enable the %s collector (default: %v)", clusterSubsystem, collector.DefaultEnabled))

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

	type metricResponse struct {
		Timestamp  string     `json:"timestamp"`
		Status     string     `json:"status"`
		Latency    perfMetric `json:"latency"`
		IOps       perfMetric `json:"iops"`
		Throughput perfMetric `json:"throughput"`
	}

	type ontapResponse struct {
		Metric metricResponse `json:"metric"`
	}

	var data *ontapResponse

	extraConfig := make(map[string]interface{}, 0)

	extraConfig["api"] = "/api/cluster/?fields=metric"
	body, err := clientAPI.Get(loginData, extraConfig)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body.(*[]byte), &data)
	if err != nil {
		return err
	}

	clusterLabels := map[string]string{"nacluster": loginData["target"].(string)}

	addMetrics(ch, "latency", namespace, clusterSubsystem, clusterLabels, &data.Metric.Latency)
	addMetrics(ch, "iops", namespace, clusterSubsystem, clusterLabels, &data.Metric.IOps)
	addMetrics(ch, "throughput", namespace, clusterSubsystem, clusterLabels, &data.Metric.Throughput)

	return nil
}

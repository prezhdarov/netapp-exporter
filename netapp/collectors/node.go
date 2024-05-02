package ontapCollectors

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/prezhdarov/prometheus-exporter/collector"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	nodeSubsystem = "node"
)

var nodeCollectorFlag = flag.Bool("node.collector", collector.DefaultEnabled, fmt.Sprintf("Enable the %s collector (default: %v)", nodeSubsystem, collector.DefaultEnabled))

type nodeCollector struct {
	logger log.Logger
}

func init() {
	collector.RegisterCollector("node", nodeCollectorFlag, NewNodeCollector)
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewNodeCollector(logger log.Logger) (collector.Collector, error) {
	return &nodeCollector{logger}, nil
}

func (c *nodeCollector) Update(ch chan<- prometheus.Metric, namespace string, clientAPI collector.ClientAPI, loginData map[string]interface{}, params map[string]string) error {

	var nodeCount = len(loginData["nodes"].(map[string]string))

	errchan := make(chan error, nodeCount)

	wg := sync.WaitGroup{}
	wg.Add(nodeCount)

	for node := range loginData["nodes"].(map[string]string) {

		go func(node string) {
			updateNode(ch, node, namespace, clientAPI, loginData, params, errchan)
			wg.Done()
		}(node)

	}

	for i := 0; i < nodeCount; i++ {
		if err := <-errchan; err != nil {
			close(errchan)
			return err
		}
	}

	close(errchan)

	wg.Wait()

	return nil

}

func updateNode(ch chan<- prometheus.Metric, node, namespace string, clientAPI collector.ClientAPI, loginData map[string]interface{}, params map[string]string, errorchan chan<- error) error {

	type metricResponse struct {
		Timestamp string  `json:"timestamp"`
		Status    string  `json:"status"`
		CPU       float64 `json:"processor_utilization"`
	}

	type failedFRU struct {
		Count float64 `json:"conunt"`
	}

	type nodeFRUs struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		State string `json:"state"`
	}
	type cpuInfo struct {
		Processor string `json:"processor"`
		Firmware  string `json:"firmware_release"`
		CPU       int    `json:"count"`
	}
	type controllerData struct {
		CPU cpuInfo `json:"cpu"`
		//	Memory      float64   `json:"memory_size"`
		Temperature string     `json:"over_temperature"`
		FailedFan   failedFRU  `json:"failed_fan"`
		FailedPSU   failedFRU  `json:"failed_power_supply"`
		FRUs        []nodeFRUs `json:"frus"`
	}

	type ontapVersion struct {
		Version string `json:"full"`
	}

	type serviceProcessor struct {
		Firmware string `json:"firmware_version"`
		State    string `json:"state"`
	}

	type nvRAM struct {
		State string `json:"battery_state"`
	}

	type ontapResponse struct {
		Model            string           `json:"model"`
		Version          ontapVersion     `json:"version"`
		Uptime           float64          `json:"uptime"`
		State            string           `json:"state"`
		Membership       string           `json:"membership"`
		StorageConfig    string           `json:"storage_configuration"`
		Controller       controllerData   `json:"controller"`
		ServiceProcessor serviceProcessor `json:"service_processor"`
		NVRAM            nvRAM            `json:"nvram"`
		Metric           metricResponse   `json:"metric"`

		//Metric metricResponse `json:"metric"`
	}

	var data *ontapResponse

	extraConfig := make(map[string]interface{}, 0)

	extraConfig["api"] = fmt.Sprintf("/api/cluster/nodes/%s?fields=model,version,uptime,state,membership,storage_configuration,controller,service_processor,nvram,metric", node)
	body, err := clientAPI.Get(loginData, extraConfig)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body.(*[]byte), &data)
	if err != nil {
		return err
	}

	nodeLabels := map[string]string{"nacluster": loginData["target"].(string), "node": loginData["nodes"].(map[string]string)[node]}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "uptime"),
			"Node uptime", nil, nodeLabels,
		), prometheus.GaugeValue, data.Uptime,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "cpu_avg"),
			"CPU Usage average in %", nil, nodeLabels,
		), prometheus.GaugeValue, data.Metric.CPU,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "fan_failed"),
			"Number of failed fans in the system", nil, nodeLabels,
		), prometheus.GaugeValue, data.Controller.FailedFan.Count,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "psu_failed"),
			"Number of failed power supplies in the system", nil, nodeLabels,
		), prometheus.GaugeValue, data.Controller.FailedPSU.Count,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "state"),
			"Node state", nil, nodeLabels,
		), prometheus.GaugeValue, nodeState[data.State],
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "sp_state"),
			"SP status", nil, nodeLabels,
		), prometheus.GaugeValue, spStatus[data.ServiceProcessor.State],
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "member_state"),
			"Node cluster membership status", nil, nodeLabels,
		), prometheus.GaugeValue, memberState[data.Membership],
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "storage_state"),
			"Node storage status", nil, nodeLabels,
		), prometheus.GaugeValue, storageState[data.StorageConfig],
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "nvram_state"),
			"Node storage status", nil, nodeLabels,
		), prometheus.GaugeValue, nvramState[data.NVRAM.State],
	)

	hardwareLabels := map[string]string{"nacluster": loginData["target"].(string), "node": loginData["nodes"].(map[string]string)[node],
		"model": data.Model, "processor": data.Controller.CPU.Processor, "cores": fmt.Sprintf("%d", data.Controller.CPU.CPU)}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "hardware_info"),
			"Node hardware information", nil, hardwareLabels,
		), prometheus.GaugeValue, float64(1),
	)

	softwareLabels := map[string]string{"nacluster": loginData["target"].(string), "node": loginData["nodes"].(map[string]string)[node],
		"version": strings.Split(data.Version.Version, ":")[0], "BIOS": data.Controller.CPU.Firmware, "SP": data.ServiceProcessor.Firmware}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeSubsystem, "software_info"),
			"Node fotware information", nil, softwareLabels,
		), prometheus.GaugeValue, float64(1),
	)

	for _, fru := range data.Controller.FRUs {
		fruLabels := map[string]string{"nacluster": loginData["target"].(string), "node": loginData["nodes"].(map[string]string)[node], "fru": fru.ID, "type": fru.Type}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, nodeSubsystem, "fru_state"),
				"Node storage status", nil, fruLabels,
			), prometheus.GaugeValue, genericState[fru.State],
		)
	}

	errorchan <- nil

	return nil
}

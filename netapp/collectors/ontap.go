package ontapCollectors

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

/*
const (

	readDescr  = "How fast the thing is giving information away"
	writeDescr = "How fast the thing is listening to others' problems"
	otherDescr = "How fast the thing is... we've no idea what"

)

var (

	genericState = map[string]float64{"error": 0, "ok": 1}
	spStatus     = map[string]float64{"offline": 0, "unknown": 1, "degraded": 2, "rebooting": 3, "updating": 4, "node_offline": 5, "sp_daemon_offline": 6, "online": 7}
	nodeState    = map[string]float64{"down": 0, "unknown": 1, "degraded": 2, "booting": 3, "taken_over": 4, "waiting_for_giveback": 5, "up": 6}
	memberState  = map[string]float64{"available": 0, "joining": 1, "member": 2}
	storageState = map[string]float64{"unknown": 0, "single_path": 1, "multi_path": 2, "mixed_path": 3, "quad_path": 4, "single_path_ha": 11, "multi_path_ha": 12, "mixed_path_ha": 13, "quad_path_ha": 14}
	nvramState   = map[string]float64{"battery_unknown": 0, "battery_not_present": 1, "battery_at_end_of_life": 2, "battery_near_end_of_life": 3, "battery_fully_discharged": 4, "battery_partially_discharged": 5, "battery_over_charged": 6, "battery_fully_charged": 7, "battery_ok": 8}

)

	type perfMetric struct {
		Read  float64 `json:"read"`
		Write float64 `json:"write"`
		Other float64 `json:"other"`
		//Total float64 `json:"total"`
	}
*/
func Load(logger log.Logger) {
	level.Info(logger).Log("msg", "Loading NetApp ONTAP collector set")
}

/*
func addMetrics(ch chan<- prometheus.Metric, name, namespace, subsystem string, labels map[string]string, metric *perfMetric) {

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterSubsystem, (name+"_read")),
			readDescr, nil, labels,
		), prometheus.GaugeValue, metric.Read,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterSubsystem, (name+"_write")),
			writeDescr, nil, labels,
		), prometheus.GaugeValue, metric.Write,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterSubsystem, (name+"_other")),
			otherDescr, nil, labels,
		), prometheus.GaugeValue, metric.Other,
	)

}
*/

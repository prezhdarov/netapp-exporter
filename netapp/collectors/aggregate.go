package ontapCollectors

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prezhdarov/prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	aggregateSubsystem string = "aggregate"
	//maxRecords         int    = 30
	expensiveRecords string = "space.block_storage.*,metric,home_node.name"
	timeout          int    = 10
	timeFormat       string = "2006-01-02T15:04:05Z"
)

var aggregateCollectorFlag = flag.Bool(fmt.Sprintf("collector.%s", aggregateSubsystem), collector.DefaultEnabled, fmt.Sprintf("Enable the %s collector (default: %v)", aggregateSubsystem, collector.DefaultEnabled))

type aggregateCollector struct {
	logger log.Logger
}

func init() {
	collector.RegisterCollector("aggregate", aggregateCollectorFlag, NewAggregateCollector)
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewAggregateCollector(logger log.Logger) (collector.Collector, error) {
	return &aggregateCollector{logger}, nil
}

func (c *aggregateCollector) Update(ch chan<- prometheus.Metric, namespace string, clientAPI collector.ClientAPI, loginData map[string]interface{}, params map[string]string) error {

	var aggregates Aggregates

	extraConfig := make(map[string]interface{}, 0)

	extraConfig["api"] = fmt.Sprintf("/api/storage/aggregates?fields=%s&return_records=true&return_timeout=%d", expensiveRecords, timeout)

	body, err := clientAPI.Get(loginData, extraConfig, c.logger)

	if err != nil {
		level.Error(c.logger).Log("Error:", err)
		return err
	}

	err = json.Unmarshal(*body.(*[]byte), &aggregates)
	if err != nil {
		level.Error(c.logger).Log("Error:", err)
		return err
	}

	for _, aggregate := range aggregates.Records {

		timestamp, err := time.Parse(timeFormat, aggregate.Metric.Timestamp)
		if err != nil {
			level.Error(c.logger).Log("Timestamp convertion failed:", err)
		}

		labels := map[string]string{"aggregate": aggregate.Name, "na_node": aggregate.Node.Name, "na_cluster": loginData["target"].(string)}

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_capacity"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.Size,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_available"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.Available,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_used"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.Used,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_inactive"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.InactiveData,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_physical_used"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.PhysicalUsed,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_used_with_snapreserve"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.UsedWithSnapReserve,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_metadata"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.Metadata,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_compact_count"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.DataCompactCount,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "block_compact_saved"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Space.Block.DataCompactSaved,
			),
		)

		// Metrics
		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "throughput_read"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.Throughput.Read,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "throughput_write"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.Throughput.Write,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "throughput_other"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.Throughput.Other,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "latency_read"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.Latency.Read,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "latency_write"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.Latency.Write,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "latency_other"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.Latency.Other,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "iops_read"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.IOps.Read,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "iops_write"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.IOps.Write,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			timestamp, prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, aggregateSubsystem, "iops_other"),
					"Aggregate capacity in bytes",
					nil, labels,
				),
				prometheus.GaugeValue, aggregate.Metric.IOps.Other,
			),
		)
	}
	return nil
}

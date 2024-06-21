package ontapCollectors

import (
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	timeFormat string = "2006-01-02T15:04:05Z"
)

func (metric *Metric) ToTimestamp(logger log.Logger) time.Time {
	timestamp, err := time.Parse(timeFormat, metric.Timestamp)
	if err != nil {
		level.Error(logger).Log("Timestamp convertion failed:", err)
	}

	return timestamp
}

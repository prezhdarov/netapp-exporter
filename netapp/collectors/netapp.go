package ontapCollectors

import (
	"fmt"
	"log/slog"
	"time"
)

const (
	timeFormat string = "2006-01-02T15:04:05Z"
)

func (metric *Metric) ToTimestamp(logger *slog.Logger) time.Time {
	timestamp, err := time.Parse(timeFormat, metric.Timestamp)
	if err != nil {
		logger.Error("Timestamp convertion failed:", fmt.Sprintf("%s", err), nil)
	}

	return timestamp
}

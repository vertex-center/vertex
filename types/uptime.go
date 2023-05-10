package types

import (
	"errors"

	"github.com/vertex-center/vertex/pkg/logger"
)

const (
	UptimeStatusOff = "off"
	UptimeStatusOn  = "on"
)

const (
	UptimeStatusFloatOff = -1.0
	UptimeStatusFloatOn  = 0.0
)

type UptimePoint struct {
	Status string `json:"status"`
}

type Uptime struct {
	Name             string        `json:"name"`
	PingURL          *string       `json:"ping_url,omitempty"`
	Current          string        `json:"current"`
	IntervalSeconds  int           `json:"interval_seconds"`
	RemainingSeconds int           `json:"remaining_seconds"`
	History          []UptimePoint `json:"history"`
}

func UptimeStatus(value float64) string {
	switch value {
	case UptimeStatusFloatOn:
		return UptimeStatusOn
	case UptimeStatusFloatOff:
		return UptimeStatusOff
	default:
		logger.Error(errors.New("incompatible uptime status")).
			AddKeyValue("value", value).
			Print()
	}
	return ""
}

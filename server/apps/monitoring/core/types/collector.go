package types

import "github.com/vertex-center/vertex/server/apps/monitoring/core/types/metric"

type Collector struct {
	IsAlive bool            `json:"is_alive"`
	Metrics []metric.Metric `json:"metrics"`
}

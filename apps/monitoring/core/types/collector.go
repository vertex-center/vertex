package types

import "github.com/vertex-center/vertex/apps/monitoring/core/types/metrics"

type Collector struct {
	IsAlive bool             `json:"is_alive"`
	Metrics []metrics.Metric `json:"metrics"`
}

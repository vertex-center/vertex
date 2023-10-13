package types

import (
	"github.com/google/uuid"
)

const (
	MetricStatusOff float64 = 0.0
	MetricStatusOn  float64 = 1.0
)

type MetricType string

const (
	MetricTypeOnOff   MetricType = "metric_type_on_off"
	MetricTypeInteger MetricType = "metric_type_number"
)

type Metric struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Type        MetricType `json:"type,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
}

type MetricsAdapterPort interface {
	// ConfigureContainer configures an container to monitor the metrics of Vertex.
	ConfigureContainer(uuid uuid.UUID) error

	// RegisterMetrics registers the metrics that can be monitored.
	RegisterMetrics(metrics []Metric)

	Set(metricID string, value interface{}, labels ...string)
	Inc(metricID string, labels ...string)
	Dec(metricID string, labels ...string)
}

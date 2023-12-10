package metrics

const (
	MetricStatusOff float64 = 0.0
	MetricStatusOn  float64 = 1.0
)

type MetricType string

const (
	MetricTypeGauge MetricType = "metric_type_gauge"
)

type Metric struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Type        MetricType `json:"type,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
}

package metric

const (
	StatusOff float64 = 0.0
	StatusOn  float64 = 1.0
)

type Type string

const (
	TypeGauge Type = "metric_type_gauge"
)

type Metric struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Type        Type     `json:"type,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

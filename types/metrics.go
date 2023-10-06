package types

import "github.com/google/uuid"

type MetricInstanceStatus int

const (
	MetricStatusOff MetricInstanceStatus = 0
	MetricStatusOn  MetricInstanceStatus = 1
)

type MetricsAdapterPort interface {
	// ConfigureInstance configures an instance to monitor the metrics of Vertex.
	ConfigureInstance(uuid uuid.UUID) error

	// UpdateInstanceStatus updates the status of an instance.
	UpdateInstanceStatus(uuid uuid.UUID, status MetricInstanceStatus)
}

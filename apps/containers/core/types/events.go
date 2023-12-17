package types

import "github.com/vertex-center/uuid"

const (
	EventNameContainersChange      = "change"
	EventNameContainerStatusChange = "status_change"
	EventNameContainerStdout       = "stdout"
	EventNameContainerStderr       = "stderr"
	EventNameContainerDownload     = "download"
)

type (
	EventContainerLog struct {
		ContainerID uuid.UUID
		Kind        string
		Message     LogLineMessage
	}

	EventContainerStatusChange struct {
		ContainerID uuid.UUID
		ServiceID   string
		Container   Container
		Name        string
		Status      string
	}

	EventContainerDeleted struct {
		ContainerID uuid.UUID
		ServiceID   string
	}

	EventContainersLoaded struct{ Count int }
	EventContainerCreated struct{}
	EventContainersChange struct{}
)

package types

import "github.com/google/uuid"

const (
	EventNameContainersChange      = "change"
	EventNameContainerStatusChange = "status_change"
	EventNameContainerStdout       = "stdout"
	EventNameContainerStderr       = "stderr"
	EventNameContainerDownload     = "download"
)

type (
	EventContainerLoaded struct {
		Container *Container
	}

	EventContainerLog struct {
		ContainerUUID uuid.UUID
		Kind          string
		Message       LogLineMessage
	}

	EventContainerStatusChange struct {
		ContainerUUID uuid.UUID
		ServiceID     string
		Container     Container
		Name          string
		Status        string
	}

	EventContainerCreated struct{}

	EventContainerDeleted struct {
		ContainerUUID uuid.UUID
		ServiceID     string
	}

	EventContainersChange struct{}

	EventContainersLoaded struct {
		Count int
	}

	EventContainersStopped struct{}
)

package types

const (
	EventNameContainersChange      = "change"
	EventNameContainerStatusChange = "status_change"
	EventNameContainerStdout       = "stdout"
	EventNameContainerStderr       = "stderr"
	EventNameContainerDownload     = "download"
)

type (
	EventContainerLog struct {
		ContainerID ContainerID
		Kind        string
		Message     LogLineMessage
	}

	EventContainerStatusChange struct {
		ContainerUUID ContainerID
		ServiceID     string
		Container     Container
		Name          string
		Status        string
	}

	EventContainerDeleted struct {
		ContainerID ContainerID
		ServiceID   string
	}

	EventContainersLoaded struct{ Count int }
	EventContainerCreated struct{}
	EventContainersChange struct{}
)

package types

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
		ContainerUUID ContainerID
		Kind          string
		Message       LogLineMessage
	}

	EventContainerStatusChange struct {
		ContainerUUID ContainerID
		ServiceID     string
		Container     Container
		Name          string
		Status        string
	}

	EventContainerDeleted struct {
		ContainerUUID ContainerID
		ServiceID     string
	}

	EventContainersLoaded struct {
		Count int
	}

	EventContainerCreated  struct{}
	EventContainersChange  struct{}
	EventContainersStopped struct{}
)

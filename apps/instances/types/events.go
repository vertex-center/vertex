package types

import "github.com/google/uuid"

const (
	EventNameInstancesChange      = "change"
	EventNameInstanceStatusChange = "status_change"
	EventNameInstanceStdout       = "stdout"
	EventNameInstanceStderr       = "stderr"
	EventNameInstanceDownload     = "download"
)

type (
	EventInstanceLoaded struct {
		Instance *Instance
	}

	EventInstanceLog struct {
		InstanceUUID uuid.UUID
		Kind         string
		Message      LogLineMessage
	}

	EventInstanceStatusChange struct {
		InstanceUUID uuid.UUID
		ServiceID    string
		Instance     Instance
		Name         string
		Status       string
	}

	EventInstanceCreated struct{}

	EventInstanceDeleted struct {
		InstanceUUID uuid.UUID
		ServiceID    string
	}

	EventInstancesChange struct{}

	EventInstancesLoaded struct {
		Count int
	}

	EventInstancesStopped struct{}
)

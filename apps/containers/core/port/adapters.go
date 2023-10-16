package port

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	types2 "github.com/vertex-center/vertex/core/types"
	"io"
)

type ContainerAdapter interface {
	Create(uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error
	GetAll() ([]uuid.UUID, error)
}

type ContainerEnvAdapter interface {
	Save(uuid uuid.UUID, env types.ContainerEnvVariables) error
	Load(uuid uuid.UUID) (types.ContainerEnvVariables, error)
}

type ContainerServiceAdapter interface {
	Save(uuid uuid.UUID, service types.Service) error
	Load(uuid uuid.UUID) (types.Service, error)
	LoadRaw(uuid uuid.UUID) (interface{}, error)
}

type ContainerSettingsAdapter interface {
	Save(uuid uuid.UUID, settings types.ContainerSettings) error
	Load(uuid uuid.UUID) (types.ContainerSettings, error)
}

type ContainerLogsAdapter interface {
	Register(uuid uuid.UUID) error
	Unregister(uuid uuid.UUID) error
	UnregisterAll() error

	Push(uuid uuid.UUID, line types.LogLine)
	Pop(uuid uuid.UUID) (types.LogLine, error)

	// LoadBuffer will load the latest logs kept in memory.
	LoadBuffer(uuid uuid.UUID) ([]types.LogLine, error)
}

type ContainerRunnerAdapter interface {
	Delete(inst *types.Container) error
	Start(inst *types.Container, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	Stop(inst *types.Container) error
	Info(inst types.Container) (map[string]any, error)
	WaitCondition(inst *types.Container, cond types2.WaitContainerCondition) error

	CheckForUpdates(inst *types.Container) error
	HasUpdateAvailable(inst types.Container) (bool, error)
	GetAllVersions(inst types.Container) ([]string, error)
}

type ServiceAdapter interface {
	// Get a service with its id. Returns ErrServiceNotFound if
	// the service was not found.
	Get(id string) (types.Service, error)

	GetScript(id string) ([]byte, error)

	// GetRaw gets a service by id, without any processing.
	// Returns ErrServiceNotFound if the service was not found.
	GetRaw(id string) (interface{}, error)

	// GetAll gets all available services.
	GetAll() []types.Service

	// Reload the adapter
	Reload() error
}

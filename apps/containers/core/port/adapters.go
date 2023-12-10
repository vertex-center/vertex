package port

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type (
	ContainerAdapter interface {
		Create(uuid uuid.UUID) error
		Delete(uuid uuid.UUID) error
		GetAll() ([]uuid.UUID, error)
	}

	ContainerEnvAdapter interface {
		Save(uuid uuid.UUID, env types.ContainerEnvVariables) error
		Load(uuid uuid.UUID) (types.ContainerEnvVariables, error)
	}

	ContainerServiceAdapter interface {
		Save(uuid uuid.UUID, service types.Service) error
		Load(uuid uuid.UUID) (types.Service, error)
		LoadRaw(uuid uuid.UUID) (interface{}, error)
	}

	ContainerSettingsAdapter interface {
		Save(uuid uuid.UUID, settings types.ContainerSettings) error
		Load(uuid uuid.UUID) (types.ContainerSettings, error)
	}

	ContainerLogsAdapter interface {
		Register(uuid uuid.UUID) error
		Unregister(uuid uuid.UUID) error
		UnregisterAll() error

		Push(uuid uuid.UUID, line types.LogLine)
		Pop(uuid uuid.UUID) (types.LogLine, error)

		// LoadBuffer will load the latest logs kept in memory.
		LoadBuffer(uuid uuid.UUID) ([]types.LogLine, error)
	}

	ContainerRunnerAdapter interface {
		DeleteContainer(ctx context.Context, inst *types.Container) error
		DeleteMounts(ctx context.Context, inst *types.Container) error
		Start(ctx context.Context, inst *types.Container, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
		Stop(ctx context.Context, inst *types.Container) error
		Info(ctx context.Context, inst types.Container) (map[string]any, error)
		WaitCondition(ctx context.Context, inst *types.Container, cond types.WaitContainerCondition) error

		CheckForUpdates(ctx context.Context, inst *types.Container) error
		HasUpdateAvailable(ctx context.Context, inst types.Container) (bool, error)
		GetAllVersions(ctx context.Context, inst types.Container) ([]string, error)
	}

	ServiceAdapter interface {
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

	DockerAdapter interface {
		ListContainers() ([]types.DockerContainer, error)
		DeleteContainer(id string) error
		CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error)
		StartContainer(id string) error
		StopContainer(id string) error
		InfoContainer(id string) (types.InfoContainerResponse, error)
		LogsStdoutContainer(id string) (io.ReadCloser, error)
		LogsStderrContainer(id string) (io.ReadCloser, error)
		WaitContainer(id string, cond types.WaitContainerCondition) error

		InfoImage(id string) (types.InfoImageResponse, error)
		PullImage(options types.PullImageOptions) (io.ReadCloser, error)
		BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error)
	}
)

package port

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type (
	ContainerAdapter interface {
		Create(uuid types.ContainerID) error
		Delete(uuid types.ContainerID) error
		GetAll() ([]types.ContainerID, error)
	}

	EnvAdapter interface {
		Save(uuid types.ContainerID, env types.ContainerEnvVariables) error
		Load(uuid types.ContainerID) (types.ContainerEnvVariables, error)
	}

	ContainerServiceAdapter interface {
		Save(uuid types.ContainerID, service types.Service) error
		Load(uuid types.ContainerID) (types.Service, error)
		LoadRaw(uuid types.ContainerID) (interface{}, error)
	}

	SettingsAdapter interface {
		Save(uuid types.ContainerID, settings types.ContainerSettings) error
		Load(uuid types.ContainerID) (types.ContainerSettings, error)
	}

	LogsAdapter interface {
		Register(uuid types.ContainerID) error
		Unregister(uuid types.ContainerID) error
		UnregisterAll() error

		Push(uuid types.ContainerID, line types.LogLine)
		Pop(uuid types.ContainerID) (types.LogLine, error)

		// LoadBuffer will load the latest logs kept in memory.
		LoadBuffer(uuid types.ContainerID) ([]types.LogLine, error)
	}

	RunnerAdapter interface {
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

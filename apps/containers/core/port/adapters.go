package port

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type (
	ContainerAdapter interface {
		GetContainer(ctx context.Context, id types.ContainerID) (*types.Container, error)
		GetContainers(ctx context.Context) (types.Containers, error)
		CreateContainer(ctx context.Context, container types.Container) error
		DeleteContainer(ctx context.Context, id types.ContainerID) error
	}

	PortAdapter interface {
		GetPorts(ctx context.Context, id types.ContainerID) (types.Ports, error)
		CreatePorts(ctx context.Context, ports types.Ports) error
		DeletePorts(ctx context.Context, id types.ContainerID) error
	}

	VolumeAdapter interface {
		GetVolumes(ctx context.Context, id types.ContainerID) (types.Volumes, error)
		CreateVolumes(ctx context.Context, volumes types.Volumes) error
		DeleteVolumes(ctx context.Context, id types.ContainerID) error
	}

	TagAdapter interface {
		CreateTags(ctx context.Context, tags types.Tags) error
		DeleteTags(ctx context.Context, id types.ContainerID) error
		GetContainerTags(ctx context.Context, id types.ContainerID) (types.Tags, error)
		GetUniqueTags(ctx context.Context) (types.Tags, error)
	}

	SysctlAdapter interface {
		GetSysctls(ctx context.Context, id types.ContainerID) (types.Sysctls, error)
		CreateSysctls(ctx context.Context, sysctls types.Sysctls) error
		DeleteSysctls(ctx context.Context, id types.ContainerID) error
	}

	EnvAdapter interface {
		GetVariable(ctx context.Context, id types.ContainerID) (types.EnvVariables, error)
		CreateVariables(ctx context.Context, env types.EnvVariables) error
		DeleteVariables(ctx context.Context, id types.ContainerID) error
	}

	CapAdapter interface {
		GetCaps(ctx context.Context, id types.ContainerID) (types.Capabilities, error)
		CreateCaps(ctx context.Context, caps types.Capabilities) error
		DeleteCaps(ctx context.Context, id types.ContainerID) error
	}

	LogsAdapter interface {
		Register(id types.ContainerID) error
		Unregister(id types.ContainerID) error
		UnregisterAll() error
		Push(id types.ContainerID, line types.LogLine)
		Pop(id types.ContainerID) (types.LogLine, error)
		LoadBuffer(id types.ContainerID) ([]types.LogLine, error) // LoadBuffer loads the latest logs kept in memory.
	}

	RunnerAdapter interface {
		DeleteContainer(ctx context.Context, c *types.Container) error
		DeleteMounts(ctx context.Context, c *types.Container) error
		Start(ctx context.Context, c *types.Container, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
		Stop(ctx context.Context, c *types.Container) error
		Info(ctx context.Context, c types.Container) (map[string]any, error)
		WaitCondition(ctx context.Context, c *types.Container, cond types.WaitContainerCondition) error
		CheckForUpdates(ctx context.Context, c *types.Container) error
		HasUpdateAvailable(ctx context.Context, c types.Container) (bool, error)
		GetAllVersions(ctx context.Context, c types.Container) ([]string, error)
	}

	ServiceAdapter interface {
		Get(id string) (types.Service, error)
		GetRaw(id string) (interface{}, error)
		GetAll() []types.Service
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

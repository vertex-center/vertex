package port

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/uuid"
)

type (
	ContainerAdapter interface {
		GetContainer(ctx context.Context, id uuid.UUID) (*types.Container, error)
		GetContainers(ctx context.Context) (types.Containers, error)
		GetContainersWithFilters(ctx context.Context, filters types.ContainerFilters) (types.Containers, error)
		CreateContainer(ctx context.Context, container types.Container) error
		UpdateContainer(ctx context.Context, container types.Container) error
		DeleteContainer(ctx context.Context, id uuid.UUID) error
		GetContainerTags(ctx context.Context, id uuid.UUID) (types.Tags, error)
		AddTag(ctx context.Context, id uuid.UUID, tagID types.TagID) error
		DeleteTags(ctx context.Context, id uuid.UUID) error
		SetStatus(ctx context.Context, id uuid.UUID, status string) error
	}

	PortAdapter interface {
		GetPorts(ctx context.Context, id uuid.UUID) (types.Ports, error)
		CreatePort(ctx context.Context, port types.Port) error
		DeletePorts(ctx context.Context, id uuid.UUID) error
	}

	VolumeAdapter interface {
		GetVolumes(ctx context.Context, id uuid.UUID) (types.Volumes, error)
		CreateVolume(ctx context.Context, vol types.Volume) error
		DeleteVolumes(ctx context.Context, id uuid.UUID) error
	}

	TagAdapter interface {
		GetTag(ctx context.Context, userID uuid.UUID, name string) (types.Tag, error)
		GetTags(ctx context.Context, userID uuid.UUID) (types.Tags, error)
		CreateTag(ctx context.Context, tag types.Tag) error
		DeleteTag(ctx context.Context, id types.TagID) error
	}

	SysctlAdapter interface {
		GetSysctls(ctx context.Context, id uuid.UUID) (types.Sysctls, error)
		CreateSysctl(ctx context.Context, sysctl types.Sysctl) error
		DeleteSysctls(ctx context.Context, id uuid.UUID) error
	}

	EnvAdapter interface {
		GetVariables(ctx context.Context, id uuid.UUID) (types.EnvVariables, error)
		CreateVariable(ctx context.Context, variable types.EnvVariable) error
		DeleteVariables(ctx context.Context, id uuid.UUID) error
		UpdateVariable(ctx context.Context, id uuid.UUID, key, value string) error
	}

	CapAdapter interface {
		GetCaps(ctx context.Context, id uuid.UUID) (types.Capabilities, error)
		CreateCap(ctx context.Context, c types.Capability) error
		DeleteCaps(ctx context.Context, id uuid.UUID) error
	}

	LogsAdapter interface {
		Register(id uuid.UUID) error
		Unregister(id uuid.UUID) error
		UnregisterAll() error
		Push(id uuid.UUID, line types.LogLine)
		Pop(id uuid.UUID) (types.LogLine, error)
		LoadBuffer(id uuid.UUID) ([]types.LogLine, error) // LoadBuffer loads the latest logs kept in memory.
		Exists(id uuid.UUID) bool
	}

	RunnerAdapter interface {
		DeleteContainer(ctx context.Context, c *types.Container) error
		DeleteMounts(ctx context.Context, c *types.Container) error
		Start(ctx context.Context, c *types.Container, ports types.Ports, volumes types.Volumes, env types.EnvVariables, caps types.Capabilities, sysctls types.Sysctls, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
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

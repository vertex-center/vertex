package port

import (
	"context"
	"io"

	vtypes "github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
)

type (
	ContainerService interface {
		Get(ctx context.Context, uuid uuid.UUID) (*types.Container, error)
		GetAll(ctx context.Context) map[uuid.UUID]*types.Container
		GetTags(ctx context.Context) []string
		Search(ctx context.Context, query types.ContainerSearchQuery) map[uuid.UUID]*types.Container
		Exists(ctx context.Context, uuid uuid.UUID) bool
		Delete(ctx context.Context, inst *types.Container) error
		StartAll(ctx context.Context)
		StopAll(ctx context.Context)
		LoadAll(ctx context.Context)
		DeleteAll(ctx context.Context)
		Install(ctx context.Context, service types.Service, method string) (*types.Container, error)
		CheckForUpdates(tx context.Context) (map[uuid.UUID]*types.Container, error)
		SetDatabases(ctx context.Context, inst *types.Container, databases map[string]uuid.UUID, options map[string]*types.SetDatabasesOptions) error
	}

	ContainerEnvService interface {
		Save(inst *types.Container, env types.ContainerEnvVariables) error
		Load(inst *types.Container) error
	}

	ContainerLogsService interface {
		GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error)
	}

	ContainerRunnerService interface {
		Install(ctx context.Context, uuid uuid.UUID, service types.Service) error
		Delete(ctx context.Context, inst *types.Container) error
		Start(ctx context.Context, inst *types.Container) error
		Stop(ctx context.Context, inst *types.Container) error
		GetDockerContainerInfo(ctx context.Context, inst types.Container) (map[string]any, error)
		GetAllVersions(ctx context.Context, inst *types.Container, useCache bool) ([]string, error)
		CheckForUpdates(ctx context.Context, inst *types.Container) error
		RecreateContainer(ctx context.Context, inst *types.Container) error
		WaitStatus(ctx context.Context, inst *types.Container, status string) error
	}

	ContainerServiceService interface {
		CheckForUpdate(inst *types.Container, latest types.Service) error
		Update(inst *types.Container, service types.Service) error
		Save(inst *types.Container, service types.Service) error
		Load(uuid uuid.UUID) (types.Service, error)
	}

	ContainerSettingsService interface {
		Save(inst *types.Container, settings types.ContainerSettings) error
		Load(inst *types.Container) error
		SetLaunchOnStartup(inst *types.Container, value bool) error
		SetDisplayName(inst *types.Container, value string) error
		SetDatabases(inst *types.Container, databases map[string]uuid.UUID) error
		SetVersion(inst *types.Container, value string) error
		SetTags(inst *types.Container, tags []string) error
	}

	MetricsService interface {
		metric.RegistryProvider
	}

	ServiceService interface {
		GetAll() []types.Service
		GetById(id string) (types.Service, error)
	}

	DockerService interface {
		ListContainers() ([]types.DockerContainer, error)
		DeleteContainer(id string) error
		CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error)
		StartContainer(id string) error
		StopContainer(id string) error
		InfoContainer(id string) (types.InfoContainerResponse, error)
		LogsStdoutContainer(id string) (io.ReadCloser, error)
		LogsStderrContainer(id string) (io.ReadCloser, error)
		WaitContainer(id string, cond types.WaitContainerCondition) error
		DeleteMounts(uuid string) error

		InfoImage(id string) (types.InfoImageResponse, error)
		PullImage(options types.PullImageOptions) (io.ReadCloser, error)
		BuildImage(options types.BuildImageOptions) (vtypes.ImageBuildResponse, error)
	}
)

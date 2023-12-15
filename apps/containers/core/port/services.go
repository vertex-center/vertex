package port

import (
	"context"
	"io"

	vtypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
)

type (
	ContainerService interface {
		Get(ctx context.Context, id types.ContainerID) (*types.Container, error)
		GetContainers(ctx context.Context) (types.Containers, error)
		Search(ctx context.Context, query types.ContainerSearchQuery) (types.Containers, error)
		Delete(ctx context.Context, id types.ContainerID) error
		UpdateContainer(ctx context.Context, id types.ContainerID, c types.Container) error
		Start(ctx context.Context, id types.ContainerID) error
		StartAll(ctx context.Context) error
		Stop(ctx context.Context, id types.ContainerID) error
		StopAll(ctx context.Context) error
		AddTag(ctx context.Context, id types.ContainerID, tagID types.TagID) error
		RecreateContainer(ctx context.Context, id types.ContainerID) error
		DeleteAll(ctx context.Context) error
		Install(ctx context.Context, serviceID string) (*types.Container, error)
		CheckForUpdates(ctx context.Context) (types.Containers, error)
		SetDatabases(ctx context.Context, c *types.Container, databases map[string]types.ContainerID, options map[string]*types.SetDatabasesOptions) error
		GetContainerEnv(ctx context.Context, id types.ContainerID) (types.EnvVariables, error)
		SaveEnv(ctx context.Context, id types.ContainerID, env types.EnvVariables) error
		GetAllVersions(ctx context.Context, id types.ContainerID, useCache bool) ([]string, error)
		GetContainerInfo(ctx context.Context, id types.ContainerID) (map[string]any, error)
		WaitStatus(ctx context.Context, id types.ContainerID, status string) error
		GetLatestLogs(id types.ContainerID) ([]types.LogLine, error)
		GetServiceByID(ctx context.Context, id string) (*types.Service, error)
		GetServices(ctx context.Context) []types.Service
	}

	MetricsService interface {
		metric.RegistryProvider
	}

	TagsService interface {
		GetTags(ctx context.Context, userID uint) (types.Tags, error)
		CreateTag(ctx context.Context, tag types.Tag) error
		DeleteTag(ctx context.Context, id types.TagID) error
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

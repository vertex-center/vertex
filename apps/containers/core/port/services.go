package port

import (
	"context"
	"io"

	vtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
)

type (
	ContainerService interface {
		Get(ctx context.Context, id uuid.UUID) (*types.Container, error)
		GetContainers(ctx context.Context) (types.Containers, error)
		GetContainersWithFilters(ctx context.Context, filters types.ContainerFilters) (types.Containers, error)
		CreateContainer(ctx context.Context, opts types.CreateContainerOptions) (*types.Container, error)
		Delete(ctx context.Context, id uuid.UUID) error
		UpdateContainer(ctx context.Context, id uuid.UUID, c types.Container) error
		Start(ctx context.Context, id uuid.UUID) error
		StartAll(ctx context.Context) error
		Stop(ctx context.Context, id uuid.UUID) error
		StopAll(ctx context.Context) error
		AddContainerTag(ctx context.Context, id uuid.UUID, tagID uuid.UUID) error
		RecreateContainer(ctx context.Context, id uuid.UUID) error
		DeleteAll(ctx context.Context) error
		CheckForUpdates(ctx context.Context) (types.Containers, error)
		SetDatabases(ctx context.Context, c *types.Container, databases map[string]uuid.UUID, options map[string]*types.SetDatabasesOptions) error
		GetContainerEnv(ctx context.Context, id uuid.UUID) (types.EnvVariables, error)
		SaveEnv(ctx context.Context, id uuid.UUID, env types.EnvVariables) error
		GetAllVersions(ctx context.Context, id uuid.UUID, useCache bool) ([]string, error)
		GetContainerInfo(ctx context.Context, id uuid.UUID) (map[string]any, error)
		WaitStatus(ctx context.Context, id uuid.UUID, status string) error
		GetLatestLogs(id uuid.UUID) ([]types.LogLine, error)
		GetTemplateByID(ctx context.Context, id string) (*types.Template, error)
		GetTemplates(ctx context.Context) []types.Template
	}

	MetricsService interface {
		metric.RegistryProvider
	}

	TagsService interface {
		GetTag(ctx context.Context, userID uuid.UUID, name string) (types.Tag, error)
		GetTags(ctx context.Context, userID uuid.UUID) (types.Tags, error)
		CreateTag(ctx context.Context, tag types.Tag) (types.Tag, error)
		DeleteTag(ctx context.Context, id uuid.UUID) error
	}

	DockerService interface {
		ListContainers() ([]types.DockerContainer, error)
		DeleteContainer(id string) error
		CreateContainer(options types.CreateDockerContainerOptions) (types.CreateContainerResponse, error)
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
		CreateVolume(name string) (volume.Volume, error)
		DeleteVolume(name string) error
	}
)

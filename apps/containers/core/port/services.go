package port

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
)

type (
	ContainerService interface {
		Get(uuid uuid.UUID) (*types.Container, error)
		GetAll() map[uuid.UUID]*types.Container
		GetTags() []string
		Search(query types.ContainerSearchQuery) map[uuid.UUID]*types.Container
		Exists(uuid uuid.UUID) bool
		Delete(inst *types.Container) error
		StartAll()
		StopAll()
		LoadAll()
		DeleteAll()
		Install(service types.Service, method string) (*types.Container, error)
		CheckForUpdates() (map[uuid.UUID]*types.Container, error)
		SetDatabases(inst *types.Container, databases map[string]uuid.UUID) error
	}

	ContainerEnvService interface {
		Save(inst *types.Container, env types.ContainerEnvVariables) error
		Load(inst *types.Container) error
	}

	ContainerLogsService interface {
		GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error)
	}

	ContainerRunnerService interface {
		Install(uuid uuid.UUID, service types.Service) error
		Delete(inst *types.Container) error
		Start(inst *types.Container) error
		Stop(inst *types.Container) error
		GetDockerContainerInfo(inst types.Container) (map[string]any, error)
		GetAllVersions(inst *types.Container, useCache bool) ([]string, error)
		CheckForUpdates(inst *types.Container) error
		RecreateContainer(inst *types.Container) error
		WaitCondition(inst *types.Container, condition vtypes.WaitContainerCondition) error
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

	MetricsService interface{}

	ServiceService interface {
		GetAll() []types.Service
		GetById(id string) (types.Service, error)
	}
)

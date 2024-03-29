package port

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
)

type (
	MockContainerService struct{ mock.Mock }
	MockTagsService      struct{ mock.Mock }
)

var (
	_ ContainerService = &MockContainerService{}
	_ TagsService      = &MockTagsService{}
)

func (m *MockContainerService) Get(ctx context.Context, id uuid.UUID) (*types.Container, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Container), args.Error(1)
}

func (m *MockContainerService) GetContainers(ctx context.Context) (types.Containers, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Containers), args.Error(1)
}

func (m *MockContainerService) GetContainersWithFilters(ctx context.Context, filters types.ContainerFilters) (types.Containers, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(types.Containers), args.Error(1)
}

func (m *MockContainerService) CreateContainer(ctx context.Context, opts types.CreateContainerOptions) (*types.Container, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Container), args.Error(1)
}

func (m *MockContainerService) Delete(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) UpdateContainer(ctx context.Context, uuid uuid.UUID, c types.Container) error {
	args := m.Called(ctx, uuid, c)
	return args.Error(0)
}

func (m *MockContainerService) Start(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) Stop(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) StartAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockContainerService) StopAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockContainerService) LoadAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockContainerService) DeleteAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockContainerService) AddContainerTag(ctx context.Context, id uuid.UUID, tagID uuid.UUID) error {
	args := m.Called(ctx, id, tagID)
	return args.Error(0)
}

func (m *MockContainerService) GetEnvs(ctx context.Context, id uuid.UUID) ([]types.EnvVariable, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]types.EnvVariable), args.Error(1)
}

func (m *MockContainerService) PatchEnv(ctx context.Context, env types.EnvVariable) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}

func (m *MockContainerService) DeleteEnv(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContainerService) CreateEnv(ctx context.Context, env types.EnvVariable) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}

func (m *MockContainerService) GetPorts(ctx context.Context, id uuid.UUID) (types.Ports, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(types.Ports), args.Error(1)
}

func (m *MockContainerService) RecreateContainer(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) CheckForUpdates(ctx context.Context) (types.Containers, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Containers), args.Error(1)
}

func (m *MockContainerService) SetDatabases(ctx context.Context, inst *types.Container, databases map[string]uuid.UUID, options map[string]*types.SetDatabasesOptions) error {
	args := m.Called(ctx, inst, databases, options)
	return args.Error(0)
}

func (m *MockContainerService) GetAllVersions(ctx context.Context, id uuid.UUID, useCache bool) ([]string, error) {
	args := m.Called(ctx, id, useCache)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockContainerService) GetContainerInfo(ctx context.Context, id uuid.UUID) (map[string]any, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockContainerService) WaitStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockContainerService) GetLatestLogs(id uuid.UUID) ([]types.LogLine, error) {
	args := m.Called(id)
	return args.Get(0).([]types.LogLine), args.Error(1)
}

func (m *MockContainerService) GetTemplateByID(ctx context.Context, id string) (*types.Template, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Template), args.Error(1)
}

func (m *MockContainerService) GetTemplates(ctx context.Context) []types.Template {
	args := m.Called(ctx)
	return args.Get(0).([]types.Template)
}

func (m *MockContainerService) ReloadContainer(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagsService) GetTag(ctx context.Context, userID uuid.UUID, name string) (types.Tag, error) {
	args := m.Called(ctx, userID, name)
	return args.Get(0).(types.Tag), args.Error(1)
}

func (m *MockTagsService) GetTags(ctx context.Context, userID uuid.UUID) (types.Tags, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(types.Tags), args.Error(1)
}

func (m *MockTagsService) CreateTag(ctx context.Context, tag types.Tag) (types.Tag, error) {
	args := m.Called(ctx, tag)
	return args.Get(0).(types.Tag), args.Error(1)
}

func (m *MockTagsService) DeleteTag(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

package port

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type (
	MockContainerService struct{ mock.Mock }
	MockTagsService      struct{ mock.Mock }
)

var (
	_ ContainerService = &MockContainerService{}
	_ TagsService      = &MockTagsService{}
)

func (m *MockContainerService) Get(ctx context.Context, id types.ContainerID) (*types.Container, error) {
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

func (m *MockContainerService) Search(ctx context.Context, query types.ContainerSearchQuery) (types.Containers, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(types.Containers), args.Error(1)
}

func (m *MockContainerService) Delete(ctx context.Context, uuid types.ContainerID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) Start(ctx context.Context, uuid types.ContainerID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) Stop(ctx context.Context, uuid types.ContainerID) error {
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

func (m *MockContainerService) AddTag(ctx context.Context, id types.ContainerID, tagID types.TagID) error {
	args := m.Called(ctx, id, tagID)
	return args.Error(0)
}

func (m *MockContainerService) RecreateContainer(ctx context.Context, uuid types.ContainerID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) Install(ctx context.Context, serviceID string) (*types.Container, error) {
	args := m.Called(ctx, serviceID)
	return args.Get(0).(*types.Container), args.Error(1)
}

func (m *MockContainerService) CheckForUpdates(ctx context.Context) (types.Containers, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Containers), args.Error(1)
}

func (m *MockContainerService) SetDatabases(ctx context.Context, inst *types.Container, databases map[string]types.ContainerID, options map[string]*types.SetDatabasesOptions) error {
	args := m.Called(ctx, inst, databases, options)
	return args.Error(0)
}

func (m *MockContainerService) SaveEnv(ctx context.Context, id types.ContainerID, env types.EnvVariables) error {
	args := m.Called(ctx, id, env)
	return args.Error(0)
}

func (m *MockContainerService) GetAllVersions(ctx context.Context, id types.ContainerID, useCache bool) ([]string, error) {
	args := m.Called(ctx, id, useCache)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockContainerService) GetContainerInfo(ctx context.Context, id types.ContainerID) (map[string]any, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockContainerService) WaitStatus(ctx context.Context, id types.ContainerID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockContainerService) GetLatestLogs(id types.ContainerID) ([]types.LogLine, error) {
	args := m.Called(id)
	return args.Get(0).([]types.LogLine), args.Error(1)
}

func (m *MockContainerService) GetServiceByID(ctx context.Context, id string) (*types.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Service), args.Error(1)
}

func (m *MockContainerService) GetServices(ctx context.Context) []types.Service {
	args := m.Called(ctx)
	return args.Get(0).([]types.Service)
}

func (m *MockTagsService) GetTags(ctx context.Context) (types.Tags, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Tags), args.Error(1)
}

func (m *MockTagsService) CreateTag(ctx context.Context, tag types.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

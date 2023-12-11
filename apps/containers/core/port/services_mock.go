package port

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type MockContainerService struct{ mock.Mock }

func (m *MockContainerService) Get(ctx context.Context, uuid types.ContainerID) (*types.Container, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Container), args.Error(1)
}

func (m *MockContainerService) GetAll(ctx context.Context) map[types.ContainerID]*types.Container {
	args := m.Called(ctx)
	return args.Get(0).(map[types.ContainerID]*types.Container)
}

func (m *MockContainerService) GetTags(ctx context.Context) []string {
	args := m.Called(ctx)
	return args.Get(0).([]string)
}

func (m *MockContainerService) Search(ctx context.Context, query types.ContainerSearchQuery) map[types.ContainerID]*types.Container {
	args := m.Called(ctx, query)
	return args.Get(0).(map[types.ContainerID]*types.Container)
}

func (m *MockContainerService) Exists(ctx context.Context, uuid types.ContainerID) bool {
	args := m.Called(ctx, uuid)
	return args.Bool(0)
}

func (m *MockContainerService) Delete(ctx context.Context, uuid types.ContainerID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockContainerService) StartAll(ctx context.Context)  { m.Called(ctx) }
func (m *MockContainerService) StopAll(ctx context.Context)   { m.Called(ctx) }
func (m *MockContainerService) LoadAll(ctx context.Context)   { m.Called(ctx) }
func (m *MockContainerService) DeleteAll(ctx context.Context) { m.Called(ctx) }

func (m *MockContainerService) Install(ctx context.Context, service types.Service, method string) (*types.Container, error) {
	args := m.Called(ctx, service, method)
	return args.Get(0).(*types.Container), args.Error(1)
}

func (m *MockContainerService) CheckForUpdates(tx context.Context) (map[types.ContainerID]*types.Container, error) {
	args := m.Called(tx)
	return args.Get(0).(map[types.ContainerID]*types.Container), args.Error(1)
}

func (m *MockContainerService) SetDatabases(ctx context.Context, inst *types.Container, databases map[string]types.ContainerID, options map[string]*types.SetDatabasesOptions) error {
	args := m.Called(ctx, inst, databases, options)
	return args.Error(0)
}

package services

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
	"testing"
)

type UpdatesServiceTestSuite struct {
	suite.Suite
	service *UpdateService

	latestBaseline types.Baseline
	updaterA       *MockUpdater
	updaterB       *MockUpdater
}

func TestUpdatesServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdatesServiceTestSuite))
}

func (suite *UpdatesServiceTestSuite) SetupTest() {
	suite.latestBaseline = types.Baseline{
		Date:         "2023-10-13",
		Version:      "v0.12.0",
		Vertex:       "v0.12.1",
		VertexClient: "v0.12.0",
	}

	suite.updaterA = &MockUpdater{}
	suite.updaterA.On("ID").Return("vertex")
	suite.updaterB = &MockUpdater{}
	suite.updaterB.On("ID").Return("vertex-client")

	updaters := []types.Updater{
		suite.updaterA,
		suite.updaterB,
	}

	adapter := &MockBaselineAdapter{}
	adapter.On("GetLatest", context.Background(), types.SettingsUpdatesChannelStable).Return(suite.latestBaseline, nil)

	suite.service = NewUpdateService(adapter, updaters)
}

func (suite *UpdatesServiceTestSuite) TestGetUpdate() {
	suite.updaterA.On("CurrentVersion").Return("v0.11.0", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.11.0", nil)

	update, err := suite.service.GetUpdate(types.SettingsUpdatesChannelStable)
	suite.NoError(err)
	suite.NotNil(update)
	suite.Equal(suite.latestBaseline, update.Baseline)
}

func (suite *UpdatesServiceTestSuite) TestGetUpdateNoUpdate() {
	suite.updaterA.On("CurrentVersion").Return("v0.12.1", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.12.0", nil)

	update, err := suite.service.GetUpdate(types.SettingsUpdatesChannelStable)
	suite.NoError(err)
	suite.Nil(update)
}

type MockBaselineAdapter struct {
	mock.Mock
}

func (a *MockBaselineAdapter) GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error) {
	args := a.Called(ctx, channel)
	return args.Get(0).(types.Baseline), args.Error(1)
}

type MockUpdater struct {
	mock.Mock
}

func (u *MockUpdater) CurrentVersion() (string, error) {
	args := u.Called()
	return args.String(0), args.Error(1)
}

func (u *MockUpdater) Install(version string) error {
	args := u.Called(version)
	return args.Error(0)
}

func (u *MockUpdater) ID() string {
	args := u.Called()
	return args.String(0)
}

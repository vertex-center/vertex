package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	types2 "github.com/vertex-center/vertex/core/types"
	"testing"
)

type UpdateServiceTestSuite struct {
	suite.Suite
	service *UpdateService

	latestBaseline types2.Baseline
	betaBaseline   types2.Baseline
	updaterA       *MockUpdater
	updaterB       *MockUpdater
	adapter        *MockBaselineAdapter
}

func TestUpdatesServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateServiceTestSuite))
}

func (suite *UpdateServiceTestSuite) SetupTest() {
	suite.latestBaseline = types2.Baseline{
		Date:         "2023-10-13",
		Version:      "v0.12.0",
		Vertex:       "v0.12.1",
		VertexClient: "v0.12.0",
	}

	suite.betaBaseline = types2.Baseline{
		Date:         "2023-10-15",
		Version:      "v0.13.0-beta",
		Vertex:       "v0.13.5-beta",
		VertexClient: "v0.13.3-beta",
	}

	suite.updaterA = &MockUpdater{}
	suite.updaterA.On("ID").Return("vertex")
	suite.updaterB = &MockUpdater{}
	suite.updaterB.On("ID").Return("vertex_client")

	updaters := []types2.Updater{
		suite.updaterA,
		suite.updaterB,
	}

	suite.adapter = &MockBaselineAdapter{}
	suite.adapter.On("GetLatest", context.Background(), types2.SettingsUpdatesChannelStable).Return(suite.latestBaseline, nil)
	suite.adapter.On("GetLatest", context.Background(), types2.SettingsUpdatesChannelBeta).Return(suite.betaBaseline, nil)

	suite.service = NewUpdateService(types2.NewVertexContext(), suite.adapter, updaters).(*UpdateService)
}

func (suite *UpdateServiceTestSuite) TestGetUpdate() {
	suite.updaterA.On("CurrentVersion").Return("v0.11.0", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.11.0", nil)

	update, err := suite.service.GetUpdate(types2.SettingsUpdatesChannelStable)
	suite.NoError(err)
	suite.NotNil(update)
	suite.Equal(suite.latestBaseline, update.Baseline)
}

func (suite *UpdateServiceTestSuite) TestGetUpdateNoUpdate() {
	suite.updaterA.On("CurrentVersion").Return("v0.12.1", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.12.0", nil)

	update, err := suite.service.GetUpdate(types2.SettingsUpdatesChannelStable)
	suite.NoError(err)
	suite.Nil(update)
}

func (suite *UpdateServiceTestSuite) TestGetUpdateBeta() {
	suite.updaterA.On("CurrentVersion").Return("v0.12.0", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.12.0", nil)

	update, err := suite.service.GetUpdate(types2.SettingsUpdatesChannelBeta)
	suite.NoError(err)
	suite.NotNil(update)
	suite.Equal(suite.betaBaseline, update.Baseline)
}

type MockBaselineAdapter struct {
	mock.Mock
}

func (a *MockBaselineAdapter) GetLatest(ctx context.Context, channel types2.SettingsUpdatesChannel) (types2.Baseline, error) {
	args := a.Called(ctx, channel)
	return args.Get(0).(types2.Baseline), args.Error(1)
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

func (u *MockUpdater) IsInstalled() bool {
	args := u.Called()
	return args.Bool(0)
}

func (u *MockUpdater) ID() string {
	args := u.Called()
	return args.String(0)
}

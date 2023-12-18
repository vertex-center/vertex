package service

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/common"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/baseline"
)

type UpdateServiceTestSuite struct {
	suite.Suite
	service *updateService

	latestBaseline baseline.Baseline
	betaBaseline   baseline.Baseline
	updaterA       *MockUpdater
	updaterB       *MockUpdater
	adapter        *port.MockBaselinesAdapter
}

func TestUpdatesServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateServiceTestSuite))
}

func (suite *UpdateServiceTestSuite) SetupTest() {
	suite.latestBaseline = baseline.Baseline{
		Date:           "2023-10-13",
		Version:        "v0.12.0",
		Vertex:         "v0.12.1",
		VertexClient:   "v0.12.0",
		VertexServices: "071bcdc8162664fb9b6c489c00277f0cce15ad87",
	}

	suite.betaBaseline = baseline.Baseline{
		Date:           "2023-10-15",
		Version:        "v0.13.0-beta",
		Vertex:         "v0.13.5-beta",
		VertexClient:   "v0.13.3-beta",
		VertexServices: "071bcdc8162664fb9b6c489c00277f0cce15ad87",
	}

	gock.Off()
	gock.New("https://bl.vx.arra.red/").
		Get("stable.json").
		Reply(http.StatusOK).
		JSON(suite.latestBaseline)
	gock.New("https://bl.vx.arra.red/").
		Get("beta.json").
		Reply(http.StatusOK).
		JSON(suite.betaBaseline)

	suite.updaterA = &MockUpdater{}
	suite.updaterA.On("ID").Return("vertex")
	suite.updaterB = &MockUpdater{}
	suite.updaterB.On("ID").Return("vertex_client")

	updaters := []types.Updater{
		suite.updaterA,
		suite.updaterB,
	}

	ctx := common.NewVertexContext(common.About{}, false)

	suite.service = NewUpdateService(apptypes.NewContext(ctx), updaters).(*updateService)
}

func (suite *UpdateServiceTestSuite) TestGetUpdate() {
	suite.updaterA.On("CurrentVersion").Return("v0.11.0", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.11.0", nil)

	update, err := suite.service.GetUpdate(baseline.ChannelStable)
	suite.Require().NoError(err)
	suite.NotNil(update)
	suite.Equal(suite.latestBaseline, update.Baseline)
}

func (suite *UpdateServiceTestSuite) TestGetUpdateNoUpdate() {
	suite.updaterA.On("CurrentVersion").Return("v0.12.1", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.12.0", nil)

	update, err := suite.service.GetUpdate(baseline.ChannelStable)
	suite.Require().NoError(err)
	suite.Nil(update)
}

func (suite *UpdateServiceTestSuite) TestGetUpdateBeta() {
	suite.updaterA.On("CurrentVersion").Return("v0.12.0", nil)
	suite.updaterB.On("CurrentVersion").Return("v0.12.0", nil)

	update, err := suite.service.GetUpdate(baseline.ChannelBeta)
	suite.Require().NoError(err)
	suite.NotNil(update)
	suite.Equal(suite.betaBaseline, update.Baseline)
}

type MockUpdater struct{ mock.Mock }

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

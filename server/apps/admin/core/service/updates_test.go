package service

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/server/common"
	apptypes "github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/common/baseline"
)

type UpdateServiceTestSuite struct {
	suite.Suite
	service *updateService

	latestBaseline baseline.Baseline
	betaBaseline   baseline.Baseline
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

	ctx := common.NewVertexContext(common.About{
		Version: "v0.12.0",
	}, false)

	suite.service = NewUpdateService(apptypes.NewContext(ctx)).(*updateService)
}

func (suite *UpdateServiceTestSuite) TestGetUpdate() {
	update, err := suite.service.GetUpdate(baseline.ChannelStable)
	suite.Require().NoError(err)
	suite.Nil(update)
}

func (suite *UpdateServiceTestSuite) TestGetUpdateBeta() {
	update, err := suite.service.GetUpdate(baseline.ChannelBeta)
	suite.Require().NoError(err)
	suite.NotNil(update)
	suite.Equal(suite.betaBaseline, update.Baseline)
}

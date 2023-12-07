package adapter

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

const baseURL = "https://bl.vx.arra.red/"

type BaselinesApiAdapterTestSuite struct {
	suite.Suite
	adapter baselinesApiAdapter
}

func TestBaselinesApiAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(BaselinesApiAdapterTestSuite))
}

func (suite *BaselinesApiAdapterTestSuite) SetupTest() {
	suite.adapter = *NewBaselinesApiAdapter().(*baselinesApiAdapter)
}

func (suite *BaselinesApiAdapterTestSuite) TestGetLatest() {
	gock.Off()
	gock.New(baseURL).
		Get("stable.json").
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"date":            "2023-10-13",
			"version":         "v0.12.0",
			"vertex":          "v0.12.1",
			"vertex_client":   "v0.12.0",
			"vertex_services": "071bcdc8162664fb9b6c489c00277f0cce15ad87",
		})

	baseline, err := suite.adapter.GetLatest(context.Background(), "stable")
	suite.Require().NoError(err)
	suite.NotEmpty(baseline)
	log.Println(baseline)
	suite.Equal("2023-10-13", baseline.Date)
	suite.Equal("v0.12.0", baseline.Version)
	suite.Equal("v0.12.1", baseline.Vertex)
	suite.Equal("v0.12.0", baseline.VertexClient)
	suite.Equal("071bcdc8162664fb9b6c489c00277f0cce15ad87", baseline.VertexServices)
}

func (suite *BaselinesApiAdapterTestSuite) TestGetLatestBeta() {
	gock.Off()
	gock.New(baseURL).
		Get("beta.json").
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"date":            "2023-10-15",
			"version":         "v0.13.0-beta",
			"vertex":          "v0.13.5-beta",
			"vertex_client":   "v0.13.3-beta",
			"vertex_services": "071bcdc8162664fb9b6c489c00277f0cce15ad87",
		})

	baseline, err := suite.adapter.GetLatest(context.Background(), "beta")
	suite.Require().NoError(err)
	suite.NotEmpty(baseline)
	suite.Equal("2023-10-15", baseline.Date)
	suite.Equal("v0.13.0-beta", baseline.Version)
	suite.Equal("v0.13.5-beta", baseline.Vertex)
	suite.Equal("v0.13.3-beta", baseline.VertexClient)
	suite.Equal("071bcdc8162664fb9b6c489c00277f0cce15ad87", baseline.VertexServices)
}

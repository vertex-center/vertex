package baseline

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

const baseURL = "https://bl.vx.arra.red/"

type BaselineTestSuite struct {
	suite.Suite
}

func TestBaselineTestSuite(t *testing.T) {
	suite.Run(t, new(BaselineTestSuite))
}

func (suite *BaselineTestSuite) TestGetVersionByID() {
	baseline := Baseline{
		Date:           "2023-10-13",
		Version:        "v0.12.0",
		Vertex:         "v0.12.1",
		VertexClient:   "v0.12.0",
		VertexServices: "071bcdc8162664fb9b6c489c00277f0cce15ad87",
	}

	vertex, err := baseline.GetVersionByID("vertex")
	suite.Require().NoError(err)
	vertexClient, err := baseline.GetVersionByID("vertex_client")
	suite.Require().NoError(err)
	vertexServices, err := baseline.GetVersionByID("vertex_services")
	suite.Require().NoError(err)

	suite.Equal("v0.12.1", vertex)
	suite.Equal("v0.12.0", vertexClient)
	suite.Equal("071bcdc8162664fb9b6c489c00277f0cce15ad87", vertexServices)
}

func (suite *BaselineTestSuite) TestFetch() {
	gock.Off()

	suite.Run("stable", func() {
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

		baseline, err := FetchLatest(context.Background(), "stable")
		suite.Require().NoError(err)
		suite.NotEmpty(baseline)
		suite.Equal("2023-10-13", baseline.Date)
		suite.Equal("v0.12.0", baseline.Version)
		suite.Equal("v0.12.1", baseline.Vertex)
		suite.Equal("v0.12.0", baseline.VertexClient)
		suite.Equal("071bcdc8162664fb9b6c489c00277f0cce15ad87", baseline.VertexServices)
	})

	suite.Run("beta", func() {
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

		baseline, err := FetchLatest(context.Background(), "beta")
		suite.Require().NoError(err)
		suite.NotEmpty(baseline)
		suite.Equal("2023-10-15", baseline.Date)
		suite.Equal("v0.13.0-beta", baseline.Version)
		suite.Equal("v0.13.5-beta", baseline.Vertex)
		suite.Equal("v0.13.3-beta", baseline.VertexClient)
		suite.Equal("071bcdc8162664fb9b6c489c00277f0cce15ad87", baseline.VertexServices)
	})
}

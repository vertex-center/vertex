package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

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

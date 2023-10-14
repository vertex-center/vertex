package adapter

import (
	"context"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

const baseURL = "https://bl.vx.quentinguidee.dev/"

type BaselinesApiAdapterTestSuite struct {
	suite.Suite
	adapter BaselinesApiAdapter
}

func TestBaselinesApiAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(BaselinesApiAdapterTestSuite))
}

func (suite *BaselinesApiAdapterTestSuite) SetupTest() {
	suite.adapter = *NewBaselinesApiAdapter().(*BaselinesApiAdapter)
}

func (suite *BaselinesApiAdapterTestSuite) TestGetLatest() {
	gock.Off()
	gock.New(baseURL).
		Get("stable.json").
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"date":           "2023-10-13",
			"version":        "v0.12.0",
			"vertex":         "v0.12.1",
			"vertexClient":   "v0.12.0",
			"vertexServices": "071bcdc8162664fb9b6c489c00277f0cce15ad87",
		})

	baseline, err := suite.adapter.GetLatest(context.Background(), "stable")
	suite.NoError(err)
	suite.NotEmpty(baseline)
}

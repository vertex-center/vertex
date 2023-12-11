package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app"
)

type ContainerServiceTestSuite struct {
	suite.Suite
	service *containerService

	containerA types.Container
	containerB types.Container
}

func TestContainerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerServiceTestSuite))
}

func (suite *ContainerServiceTestSuite) SetupTest() {
	suite.service = NewContainerService(ContainerServiceParams{
		Ctx: app.NewContext(common.NewVertexContext(common.About{}, false)),
	}).(*containerService)

	suite.containerA = types.Container{
		UUID: types.NewContainerID(),
		Service: types.Service{
			Name: "service-a",
		},
		ContainerSettings: types.ContainerSettings{
			Tags: []string{"Global Tag", "Service A Tag 0", "Service A Tag 1"},
		},
	}

	suite.containerB = types.Container{
		UUID: types.NewContainerID(),
		Service: types.Service{
			Name: "service-b",
			Features: &types.Features{
				Databases: &[]types.DatabaseFeature{
					{
						Type: "postgres",
					},
				},
			},
		},
		ContainerSettings: types.ContainerSettings{
			Tags: []string{"Global Tag"},
		},
	}

	suite.service.containers = map[types.ContainerID]*types.Container{
		suite.containerA.UUID: &suite.containerA,
		suite.containerB.UUID: &suite.containerB,
	}
}

func (suite *ContainerServiceTestSuite) TestSearch() {
	tests := []struct {
		query    types.ContainerSearchQuery
		expected []types.ContainerID
	}{
		// Empty query
		{
			types.ContainerSearchQuery{},
			[]types.ContainerID{
				suite.containerA.UUID,
				suite.containerB.UUID,
			},
		},
		// Tags
		{
			types.ContainerSearchQuery{
				Tags: &[]string{},
			},
			[]types.ContainerID{},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Invalid Tag"},
			},
			[]types.ContainerID{},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Service A Tag 1"},
			},
			[]types.ContainerID{suite.containerA.UUID},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Service A Tag 0", "Service A Tag 1"},
			},
			[]types.ContainerID{suite.containerA.UUID},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Global Tag"},
			},
			[]types.ContainerID{
				suite.containerA.UUID,
				suite.containerB.UUID,
			},
		},
		// Features
		{
			types.ContainerSearchQuery{
				Features: &[]string{"invalid-feature"},
			},
			[]types.ContainerID{},
		},
		{
			types.ContainerSearchQuery{
				Features: &[]string{"postgres"},
			},
			[]types.ContainerID{suite.containerB.UUID},
		},
		// Multiple
		{
			types.ContainerSearchQuery{
				Tags:     &[]string{"Global Tag"},
				Features: &[]string{"postgres"},
			},
			[]types.ContainerID{suite.containerB.UUID},
		},
	}

	for _, t := range tests {
		results := suite.service.Search(context.Background(), t.query)

		suite.Len(results, len(t.expected))

		var resultUUIDs []types.ContainerID
		for id := range results {
			resultUUIDs = append(resultUUIDs, id)
		}

		for _, expected := range t.expected {
			suite.Contains(resultUUIDs, expected)
		}
	}
}

func (suite *ContainerServiceTestSuite) TestGetTags() {
	tags := suite.service.GetTags(context.Background())

	suite.Len(tags, 3)
	suite.Contains(tags, "Global Tag")
	suite.Contains(tags, "Service A Tag 0")
	suite.Contains(tags, "Service A Tag 1")
}

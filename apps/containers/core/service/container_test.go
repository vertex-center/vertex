package service

import (
	"testing"

	"github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ContainerServiceTestSuite struct {
	suite.Suite
	service *ContainerService

	containerA types.Container
	containerB types.Container
}

func TestContainerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerServiceTestSuite))
}

func (suite *ContainerServiceTestSuite) SetupTest() {
	suite.service = NewContainerService(ContainerServiceParams{
		Ctx: app.NewContext(vtypes.NewVertexContext()),
	}).(*ContainerService)

	suite.containerA = types.Container{
		UUID: uuid.New(),
		Service: types.Service{
			Name: "service-a",
		},
		ContainerSettings: types.ContainerSettings{
			Tags: []string{"Global Tag", "Service A Tag 0", "Service A Tag 1"},
		},
	}

	suite.containerB = types.Container{
		UUID: uuid.New(),
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

	suite.service.containers = map[uuid.UUID]*types.Container{
		suite.containerA.UUID: &suite.containerA,
		suite.containerB.UUID: &suite.containerB,
	}
}

func (suite *ContainerServiceTestSuite) TestSearch() {
	tests := []struct {
		query    types.ContainerSearchQuery
		expected []uuid.UUID
	}{
		// Empty query
		{
			types.ContainerSearchQuery{},
			[]uuid.UUID{
				suite.containerA.UUID,
				suite.containerB.UUID,
			},
		},
		// Tags
		{
			types.ContainerSearchQuery{
				Tags: &[]string{},
			},
			[]uuid.UUID{},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Invalid Tag"},
			},
			[]uuid.UUID{},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Service A Tag 1"},
			},
			[]uuid.UUID{suite.containerA.UUID},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Service A Tag 0", "Service A Tag 1"},
			},
			[]uuid.UUID{suite.containerA.UUID},
		},
		{
			types.ContainerSearchQuery{
				Tags: &[]string{"Global Tag"},
			},
			[]uuid.UUID{
				suite.containerA.UUID,
				suite.containerB.UUID,
			},
		},
		// Features
		{
			types.ContainerSearchQuery{
				Features: &[]string{"invalid-feature"},
			},
			[]uuid.UUID{},
		},
		{
			types.ContainerSearchQuery{
				Features: &[]string{"postgres"},
			},
			[]uuid.UUID{suite.containerB.UUID},
		},
		// Multiple
		{
			types.ContainerSearchQuery{
				Tags:     &[]string{"Global Tag"},
				Features: &[]string{"postgres"},
			},
			[]uuid.UUID{suite.containerB.UUID},
		},
	}

	for _, t := range tests {
		results := suite.service.Search(t.query)

		suite.Len(results, len(t.expected))

		var resultUUIDs []uuid.UUID
		for id := range results {
			resultUUIDs = append(resultUUIDs, id)
		}

		for _, expected := range t.expected {
			suite.Contains(resultUUIDs, expected)
		}
	}
}

func (suite *ContainerServiceTestSuite) TestGetTags() {
	tags := suite.service.GetTags()

	suite.Len(tags, 3)
	suite.Contains(tags, "Global Tag")
	suite.Contains(tags, "Service A Tag 0")
	suite.Contains(tags, "Service A Tag 1")
}

package service

import (
	"testing"

	types2 "github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ContainerServiceTestSuite struct {
	suite.Suite
	service *ContainerService

	containerA types2.Container
	containerB types2.Container
}

func TestContainerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerServiceTestSuite))
}

func (suite *ContainerServiceTestSuite) SetupTest() {
	suite.service = NewContainerService(ContainerServiceParams{
		Ctx: app.NewContext(vtypes.NewVertexContext()),
	}).(*ContainerService)

	suite.containerA = types2.Container{
		UUID: uuid.New(),
		Service: types2.Service{
			Name: "service-a",
		},
		ContainerSettings: types2.ContainerSettings{
			Tags: []string{"Global Tag", "Service A Tag 0", "Service A Tag 1"},
		},
	}

	suite.containerB = types2.Container{
		UUID: uuid.New(),
		Service: types2.Service{
			Name: "service-b",
			Features: &types2.Features{
				Databases: &[]types2.DatabaseFeature{
					{
						Type: "postgres",
					},
				},
			},
		},
		ContainerSettings: types2.ContainerSettings{
			Tags: []string{"Global Tag"},
		},
	}

	suite.service.containers = map[uuid.UUID]*types2.Container{
		suite.containerA.UUID: &suite.containerA,
		suite.containerB.UUID: &suite.containerB,
	}
}

func (suite *ContainerServiceTestSuite) TestSearch() {
	tests := []struct {
		query    types2.ContainerSearchQuery
		expected []uuid.UUID
	}{
		// Empty query
		{
			types2.ContainerSearchQuery{},
			[]uuid.UUID{
				suite.containerA.UUID,
				suite.containerB.UUID,
			},
		},
		// Tags
		{
			types2.ContainerSearchQuery{
				Tags: &[]string{},
			},
			[]uuid.UUID{},
		},
		{
			types2.ContainerSearchQuery{
				Tags: &[]string{"Invalid Tag"},
			},
			[]uuid.UUID{},
		},
		{
			types2.ContainerSearchQuery{
				Tags: &[]string{"Service A Tag 1"},
			},
			[]uuid.UUID{suite.containerA.UUID},
		},
		{
			types2.ContainerSearchQuery{
				Tags: &[]string{"Service A Tag 0", "Service A Tag 1"},
			},
			[]uuid.UUID{suite.containerA.UUID},
		},
		{
			types2.ContainerSearchQuery{
				Tags: &[]string{"Global Tag"},
			},
			[]uuid.UUID{
				suite.containerA.UUID,
				suite.containerB.UUID,
			},
		},
		// Features
		{
			types2.ContainerSearchQuery{
				Features: &[]string{"invalid-feature"},
			},
			[]uuid.UUID{},
		},
		{
			types2.ContainerSearchQuery{
				Features: &[]string{"postgres"},
			},
			[]uuid.UUID{suite.containerB.UUID},
		},
		// Multiple
		{
			types2.ContainerSearchQuery{
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

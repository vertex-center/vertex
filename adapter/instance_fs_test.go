package adapter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
)

const (
	PathInstances = "tests/instances"
)

type InstanceFSAdapterTestSuite struct {
	suite.Suite

	adapter InstanceFSAdapter

	instanceA types.Instance
	instanceB types.Instance
}

func TestInstanceFSAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceFSAdapterTestSuite))
}

func (suite *InstanceFSAdapterTestSuite) SetupSuite() {
	suite.adapter = *NewInstanceFSAdapter(&InstanceFSAdapterParams{
		instancesPath: PathInstances,
	}).(*InstanceFSAdapter)

	suite.instanceA = types.Instance{
		UUID: uuid.New(),
		Service: types.Service{
			Features: &types.Features{
				Databases: &[]types.DatabaseFeature{
					{Type: "postgres"},
				},
			},
		},
	}

	suite.instanceB = types.Instance{
		UUID: uuid.New(),
		Service: types.Service{
			Features: &types.Features{
				Databases: &[]types.DatabaseFeature{
					{Type: "redis"},
				},
			},
		},
	}

	suite.adapter.instances = map[uuid.UUID]*types.Instance{
		suite.instanceA.UUID: &suite.instanceA,
		suite.instanceB.UUID: &suite.instanceB,
	}
}

func (suite *InstanceFSAdapterTestSuite) TestGet() {
	instanceB, err := suite.adapter.Get(suite.instanceB.UUID)
	suite.NoError(err)
	suite.Equal(suite.instanceB.UUID, instanceB.UUID)
}

func (suite *InstanceFSAdapterTestSuite) TestGetAll() {
	instances := suite.adapter.GetAll()
	suite.Equal(2, len(instances))
}

func (suite *InstanceFSAdapterTestSuite) TestSearch() {
	instances := suite.adapter.Search(types.InstanceQuery{})
	suite.Equal(len(suite.adapter.instances), len(instances))

	instances = suite.adapter.Search(types.InstanceQuery{
		Features: []string{"redis"},
	})
	suite.Equal(1, len(instances))
	suite.Equal(suite.instanceB.UUID, instances[suite.instanceB.UUID].UUID)
}

func (suite *InstanceFSAdapterTestSuite) TestGetPath() {
	p := suite.adapter.GetPath(suite.instanceA.UUID)
	suite.Equal(fmt.Sprintf("tests/instances/%s", suite.instanceA.UUID), p)
}

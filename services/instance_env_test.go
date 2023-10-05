package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
)

type InstanceEnvServiceTestSuite struct {
	suite.Suite

	service InstanceEnvService
	adapter MockInstanceEnvAdapter
}

func TestInstanceEnvServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceEnvServiceTestSuite))
}

func (suite *InstanceEnvServiceTestSuite) SetupSuite() {
	suite.adapter = MockInstanceEnvAdapter{}
	suite.service = NewInstanceEnvService(&suite.adapter)
}

func (suite *InstanceEnvServiceTestSuite) TestSave() {
	suite.adapter.On("Save", mock.Anything, mock.Anything).Return(nil)

	inst := &types.Instance{}
	env := types.InstanceEnvVariables{"a": "b"}
	err := suite.service.Save(inst, env)

	suite.NoError(err)
	suite.Equal(env, inst.Env)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *InstanceEnvServiceTestSuite) TestLoad() {
	suite.adapter.On("Load", mock.Anything).Return(types.InstanceEnvVariables{}, nil)

	inst := &types.Instance{}
	err := suite.service.Load(inst)

	suite.NoError(err)
	suite.Equal(inst.Env, types.InstanceEnvVariables{"a": "b"})
	suite.adapter.AssertExpectations(suite.T())
}

type MockInstanceEnvAdapter struct {
	mock.Mock
}

func (m *MockInstanceEnvAdapter) Save(uuid uuid.UUID, env types.InstanceEnvVariables) error {
	args := m.Called(uuid, env)
	return args.Error(0)
}

func (m *MockInstanceEnvAdapter) Load(uuid uuid.UUID) (types.InstanceEnvVariables, error) {
	args := m.Called(uuid)
	return types.InstanceEnvVariables{"a": "b"}, args.Error(1)
}

package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type ContainerEnvServiceTestSuite struct {
	suite.Suite

	service *envService
	adapter MockContainerEnvAdapter
}

func TestContainerEnvServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerEnvServiceTestSuite))
}

func (suite *ContainerEnvServiceTestSuite) SetupSuite() {
	suite.adapter = MockContainerEnvAdapter{}
	suite.service = NewEnvService(&suite.adapter).(*envService)
}

func (suite *ContainerEnvServiceTestSuite) TestSave() {
	suite.adapter.On("Save", mock.Anything, mock.Anything).Return(nil)

	inst := &types.Container{}
	env := types.ContainerEnvVariables{"a": "b"}
	err := suite.service.Save(inst, env)

	suite.Require().NoError(err)
	suite.Equal(env, inst.Env)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *ContainerEnvServiceTestSuite) TestLoad() {
	suite.adapter.On("Load", mock.Anything).Return(types.ContainerEnvVariables{}, nil)

	inst := &types.Container{}
	err := suite.service.Load(inst)

	suite.Require().NoError(err)
	suite.Equal(types.ContainerEnvVariables{"a": "b"}, inst.Env)
	suite.adapter.AssertExpectations(suite.T())
}

type MockContainerEnvAdapter struct{ mock.Mock }

func (m *MockContainerEnvAdapter) Save(uuid types.ContainerID, env types.ContainerEnvVariables) error {
	args := m.Called(uuid, env)
	return args.Error(0)
}

func (m *MockContainerEnvAdapter) Load(uuid types.ContainerID) (types.ContainerEnvVariables, error) {
	args := m.Called(uuid)
	return types.ContainerEnvVariables{"a": "b"}, args.Error(1)
}

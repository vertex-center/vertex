package service

import (
	"testing"

	types2 "github.com/vertex-center/vertex/apps/containers/core/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ContainerEnvServiceTestSuite struct {
	suite.Suite

	service *ContainerEnvService
	adapter MockContainerEnvAdapter
}

func TestContainerEnvServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerEnvServiceTestSuite))
}

func (suite *ContainerEnvServiceTestSuite) SetupSuite() {
	suite.adapter = MockContainerEnvAdapter{}
	suite.service = NewContainerEnvService(&suite.adapter).(*ContainerEnvService)
}

func (suite *ContainerEnvServiceTestSuite) TestSave() {
	suite.adapter.On("Save", mock.Anything, mock.Anything).Return(nil)

	inst := &types2.Container{}
	env := types2.ContainerEnvVariables{"a": "b"}
	err := suite.service.Save(inst, env)

	suite.NoError(err)
	suite.Equal(env, inst.Env)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *ContainerEnvServiceTestSuite) TestLoad() {
	suite.adapter.On("Load", mock.Anything).Return(types2.ContainerEnvVariables{}, nil)

	inst := &types2.Container{}
	err := suite.service.Load(inst)

	suite.NoError(err)
	suite.Equal(inst.Env, types2.ContainerEnvVariables{"a": "b"})
	suite.adapter.AssertExpectations(suite.T())
}

type MockContainerEnvAdapter struct {
	mock.Mock
}

func (m *MockContainerEnvAdapter) Save(uuid uuid.UUID, env types2.ContainerEnvVariables) error {
	args := m.Called(uuid, env)
	return args.Error(0)
}

func (m *MockContainerEnvAdapter) Load(uuid uuid.UUID) (types2.ContainerEnvVariables, error) {
	args := m.Called(uuid)
	return types2.ContainerEnvVariables{"a": "b"}, args.Error(1)
}

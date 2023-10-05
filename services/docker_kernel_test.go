package services

import (
	"io"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
)

type DockerKernelServiceTestSuite struct {
	suite.Suite

	service DockerKernelService
	adapter MockDockerAdapter
}

func TestDockerKernelServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DockerKernelServiceTestSuite))
}

func (suite *DockerKernelServiceTestSuite) SetupSuite() {
	suite.adapter = MockDockerAdapter{}
	suite.service = NewDockerKernelService(&suite.adapter)
}

func (suite *DockerKernelServiceTestSuite) TestListContainers() {
	suite.adapter.On("ListContainers").Return([]types.Container{}, nil)

	containers, err := suite.service.ListContainers()

	suite.NoError(err)
	suite.Equal([]types.Container{}, containers)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestDeleteContainer() {
	suite.adapter.On("DeleteContainer", mock.Anything).Return(nil)

	err := suite.service.DeleteContainer("")

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestCreateContainer() {
	suite.adapter.On("CreateContainer", mock.Anything).Return(types.CreateContainerResponse{}, nil)

	container, err := suite.service.CreateContainer(types.CreateContainerOptions{})

	suite.NoError(err)
	suite.Equal(types.CreateContainerResponse{}, container)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestStartContainer() {
	suite.adapter.On("StartContainer", mock.Anything).Return(nil)

	err := suite.service.StartContainer("")

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestStopContainer() {
	suite.adapter.On("StopContainer", mock.Anything).Return(nil)

	err := suite.service.StopContainer("")

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestInfoContainer() {
	suite.adapter.On("InfoContainer", mock.Anything).Return(types.InfoContainerResponse{}, nil)

	info, err := suite.service.InfoContainer("")

	suite.NoError(err)
	suite.Equal(types.InfoContainerResponse{}, info)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestLogsStdoutContainer() {
	suite.adapter.On("LogsStdoutContainer", mock.Anything).Return(nil, nil)

	stdout, err := suite.service.LogsStdoutContainer("")

	suite.NoError(err)
	suite.Nil(stdout)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestLogsStderrContainer() {
	suite.adapter.On("LogsStderrContainer", mock.Anything).Return(nil, nil)

	stderr, err := suite.service.LogsStderrContainer("")

	suite.NoError(err)
	suite.Nil(stderr)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestWaitContainer() {
	suite.adapter.On("WaitContainer", mock.Anything, mock.Anything).Return(nil)

	err := suite.service.WaitContainer("", types.WaitContainerCondition(container.WaitConditionNotRunning))

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestInfoImage() {
	suite.adapter.On("InfoImage", mock.Anything).Return(types.InfoImageResponse{}, nil)

	info, err := suite.service.InfoImage("")

	suite.NoError(err)
	suite.Equal(types.InfoImageResponse{}, info)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestPullImage() {
	suite.adapter.On("PullImage", mock.Anything).Return(nil, nil)

	image, err := suite.service.PullImage(types.PullImageOptions{})

	suite.NoError(err)
	suite.Nil(image)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *DockerKernelServiceTestSuite) TestBuildImage() {
	suite.adapter.On("BuildImage", mock.Anything).Return(dockertypes.ImageBuildResponse{}, nil)

	image, err := suite.service.BuildImage(types.BuildImageOptions{})

	suite.NoError(err)
	suite.Equal(dockertypes.ImageBuildResponse{}, image)
	suite.adapter.AssertExpectations(suite.T())
}

type MockDockerAdapter struct {
	mock.Mock
}

func (m *MockDockerAdapter) ListContainers() ([]types.Container, error) {
	args := m.Called()
	return args.Get(0).([]types.Container), args.Error(1)
}

func (m *MockDockerAdapter) DeleteContainer(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDockerAdapter) CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error) {
	args := m.Called(options)
	return args.Get(0).(types.CreateContainerResponse), args.Error(1)
}

func (m *MockDockerAdapter) StartContainer(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDockerAdapter) StopContainer(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDockerAdapter) InfoContainer(id string) (types.InfoContainerResponse, error) {
	args := m.Called(id)
	return args.Get(0).(types.InfoContainerResponse), args.Error(1)
}

func (m *MockDockerAdapter) LogsStdoutContainer(id string) (io.ReadCloser, error) {
	args := m.Called(id)
	return nil, args.Error(1)
}

func (m *MockDockerAdapter) LogsStderrContainer(id string) (io.ReadCloser, error) {
	args := m.Called(id)
	return nil, args.Error(1)
}

func (m *MockDockerAdapter) WaitContainer(id string, cond types.WaitContainerCondition) error {
	args := m.Called(id, cond)
	return args.Error(0)
}

func (m *MockDockerAdapter) InfoImage(id string) (types.InfoImageResponse, error) {
	args := m.Called(id)
	return args.Get(0).(types.InfoImageResponse), args.Error(1)
}

func (m *MockDockerAdapter) PullImage(options types.PullImageOptions) (io.ReadCloser, error) {
	args := m.Called(options)
	return nil, args.Error(1)
}

func (m *MockDockerAdapter) BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error) {
	args := m.Called(options)
	return args.Get(0).(dockertypes.ImageBuildResponse), args.Error(1)
}

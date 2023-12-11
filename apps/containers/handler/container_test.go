package handler

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router/routertest"
)

type ContainerHandlerTestSuite struct {
	suite.Suite

	service port.MockContainerService
	handler *containerHandler

	testContainer types.Container
	opts          routertest.RequestOptions
}

func TestContainerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerHandlerTestSuite))
}

func (suite *ContainerHandlerTestSuite) SetupSubTest() {
	suite.service = port.MockContainerService{}
	suite.handler = NewContainerHandler(ContainerHandlerParams{
		ContainerService: &suite.service,
	}).(*containerHandler)
	suite.testContainer = types.Container{
		UUID: uuid.New(),
	}
	suite.opts = routertest.RequestOptions{
		Params: map[string]string{
			"container_uuid": suite.testContainer.UUID.String(),
		},
	}
}

func (suite *ContainerHandlerTestSuite) TestGet() {
	suite.Run("Success", func() {
		suite.service.On("Get", mock.Anything, suite.testContainer.UUID).Return(&suite.testContainer, nil)

		res := routertest.Request("GET", suite.handler.Get(), suite.opts)

		suite.Equal(200, res.Code)
		suite.service.AssertExpectations(suite.T())
	})

	suite.Run("NotFound", func() {
		suite.service.On("Get", mock.Anything, suite.testContainer.UUID).Return(nil, types.ErrContainerNotFound)

		res := routertest.Request("GET", suite.handler.Get(), suite.opts)

		suite.Equal(404, res.Code)
		suite.service.AssertExpectations(suite.T())
	})
}

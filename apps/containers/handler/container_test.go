package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/uuid"
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
	ctx := common.NewVertexContext(common.About{}, false)
	appCtx := app.NewContext(ctx)
	suite.handler = NewContainerHandler(appCtx, &suite.service).(*containerHandler)
	suite.testContainer = types.Container{
		ID: uuid.New(),
	}
	suite.opts = routertest.RequestOptions{
		Params: map[string]string{
			"container_id": suite.testContainer.ID.String(),
		},
	}
}

func (suite *ContainerHandlerTestSuite) TestGet() {
	suite.Run("OK", func() {
		suite.service.On("Get", mock.Anything, suite.testContainer.ID).Return(&suite.testContainer, nil)

		res := routertest.Request("GET", suite.handler.Get(), suite.opts)

		suite.Equal(200, res.Code)
		suite.service.AssertExpectations(suite.T())
	})

	suite.Run("Not Found", func() {
		suite.service.On("Get", mock.Anything, suite.testContainer.ID).Return(nil, types.ErrContainerNotFound)

		res := routertest.Request("GET", suite.handler.Get(), suite.opts)

		suite.Equal(404, res.Code)
		suite.JSONEq(`{"error":"container not found"}`, res.Body.String())
		suite.service.AssertExpectations(suite.T())
	})
}

func (suite *ContainerHandlerTestSuite) TestDelete() {
	suite.Run("OK", func() {
		suite.service.On("Delete", mock.Anything, suite.testContainer.ID).Return(nil)

		res := routertest.Request("DELETE", suite.handler.Delete(), suite.opts)

		suite.Equal(204, res.Code)
		suite.service.AssertExpectations(suite.T())
	})

	suite.Run("Not Found", func() {
		suite.service.On("Delete", mock.Anything, suite.testContainer.ID).Return(types.ErrContainerNotFound)

		res := routertest.Request("DELETE", suite.handler.Delete(), suite.opts)

		suite.Equal(404, res.Code)
		suite.JSONEq(`{"error":"container not found"}`, res.Body.String())
		suite.service.AssertExpectations(suite.T())
	})

	suite.Run("Error", func() {
		suite.service.On("Delete", mock.Anything, suite.testContainer.ID).Return(errors.New("error"))

		res := routertest.Request("DELETE", suite.handler.Delete(), suite.opts)

		suite.Equal(500, res.Code)
		suite.JSONEq(`{"error":"error"}`, res.Body.String())
		suite.service.AssertExpectations(suite.T())
	})
}

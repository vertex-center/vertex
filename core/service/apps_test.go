package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

type AppsServiceTestSuite struct {
	suite.Suite
	service *AppsService
	app     *MockApp
}

func TestAppsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AppsServiceTestSuite))
}

func (suite *AppsServiceTestSuite) SetupTest() {
	ctx := types.NewVertexContext(types.About{})
	suite.app = &MockApp{}
	suite.service = NewAppsService(ctx, false, router.New(), []app.Interface{
		suite.app,
	}).(*AppsService)
}

func (suite *AppsServiceTestSuite) TestStartApps() {
	suite.app.On("Load", mock.Anything).Return()
	suite.app.On("Meta").Return(app.Meta{ID: "test"})
	suite.app.On("Initialize", mock.Anything).Return(nil)
	suite.service.StartApps()
	suite.app.AssertExpectations(suite.T())
}

type MockApp struct {
	mock.Mock
}

func (m *MockApp) Load(ctx *app.Context) {
	m.Called(ctx)
}

func (m *MockApp) Meta() app.Meta {
	args := m.Called()
	return args.Get(0).(app.Meta)
}

func (m *MockApp) Initialize(r *router.Group) error {
	args := m.Called(r)
	return args.Error(0)
}

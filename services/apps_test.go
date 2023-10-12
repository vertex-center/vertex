package services

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/app"
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
	ctx := types.NewVertexContext()
	suite.app = &MockApp{}
	suite.service = NewAppsService(ctx, router.New(), []app.Interface{
		suite.app,
	})
}

func (suite *AppsServiceTestSuite) TestStartApps() {
	a := app.New(suite.service.ctx)

	suite.app.On("Initialize", a).Return(nil)
	suite.service.StartApps()
	suite.app.AssertExpectations(suite.T())
}

type MockApp struct {
	mock.Mock
}

func (m *MockApp) Initialize(app *app.App) error {
	args := m.Called(app)
	return args.Error(0)
}

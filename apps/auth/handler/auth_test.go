package handler

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/routertest"
)

type AuthHandlerTestSuite struct {
	suite.Suite

	service     port.MockAuthService
	handler     *AuthHandler
	testSession types.Session
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.service = port.MockAuthService{}
	suite.handler = NewAuthHandler(&suite.service).(*AuthHandler)
	suite.testSession = types.Session{
		Token: "test_token",
	}
}

func (suite *AuthHandlerTestSuite) TestLogin() {
	suite.service.On("Login", "test_login", "test_password").Return(suite.testSession, nil)

	auth := base64.StdEncoding.EncodeToString([]byte("test_login:test_password"))

	res := routertest.Request("POST", suite.handler.Login, routertest.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Basic " + auth,
		},
	})

	suite.Equal(200, res.Code)
	suite.JSONEq(routertest.ToJSON(suite.testSession), res.Body.String())
	suite.service.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestLoginInvalidCredentials() {
	suite.service.On("Login", "test_login", "invalid_password").Return(types.Session{}, types.ErrLoginFailed)

	auth := base64.StdEncoding.EncodeToString([]byte("test_login:invalid_password"))

	res := routertest.Request("POST", suite.handler.Login, routertest.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Basic " + auth,
		},
	})

	suite.Equal(500, res.Code)
	suite.service.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestRegister() {

	suite.service.On("Register", "test_login", "test_password").Return(suite.testSession, nil)

	auth := base64.StdEncoding.EncodeToString([]byte("test_login:test_password"))

	res := routertest.Request("POST", suite.handler.Register, routertest.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Basic " + auth,
		},
	})

	suite.Equal(200, res.Code)
	suite.JSONEq(routertest.ToJSON(suite.testSession), res.Body.String())
	suite.service.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestLogout() {
	suite.service.On("Logout", "test_token").Return(nil)

	handle := func(c *router.Context) {
		c.Set("token", "test_token")
		suite.handler.Logout(c)
	}

	res := routertest.Request("POST", handle, routertest.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Bearer test_token",
		},
	})

	suite.Equal(204, res.Code)
	suite.service.AssertExpectations(suite.T())
}

package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/config"
)

type AuthServiceTestSuite struct {
	suite.Suite

	service  *AuthService
	adapter  port.MockAuthAdapter
	testCred types.CredentialsArgon2id
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.adapter = port.MockAuthAdapter{}
	suite.service = NewAuthService(&suite.adapter).(*AuthService)
	suite.testCred = types.CredentialsArgon2id{
		Login:       "test_login",
		Hash:        "N6WEEZ++Gh54U9jqEwSmFAWz9Ls+8iyHar4mOU7M71Y=",
		Type:        "argon2id",
		Salt:        "vertex",
		Iterations:  3,
		Memory:      12 * 1024,
		Parallelism: 4,
		KeyLen:      32,
		Users: []*types.User{
			{
				ID:       1,
				Username: "test_username",
			},
		},
	}
}

func (suite *AuthServiceTestSuite) TestLogin() {
	suite.adapter.On("GetCredentials", "test_login").Return([]types.CredentialsArgon2id{suite.testCred}, nil)
	suite.adapter.On("SaveToken", mock.Anything).Return(nil)

	token, err := suite.service.Login("test_login", "test_password")
	suite.Require().NoError(err)
	suite.Equal("test_username", token.User.Username)
	suite.NotEmpty(token.Token)
	suite.Len(token.Token, 44)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestLoginInvalidLogin() {
	suite.adapter.On("GetCredentials", "invalid_login").Return([]types.CredentialsArgon2id{}, nil)

	_, err := suite.service.Login("invalid_login", "test_password")
	suite.Require().ErrorIs(err, types.ErrLoginFailed)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestLoginInvalidPassword() {
	suite.adapter.On("GetCredentials", "test_login").Return([]types.CredentialsArgon2id{suite.testCred}, nil)

	_, err := suite.service.Login("test_login", "invalid_password")
	suite.Require().ErrorIs(err, types.ErrLoginFailed)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestRegister() {
	suite.adapter.On("GetCredentials", "test_login").Return([]types.CredentialsArgon2id{suite.testCred}, nil)
	suite.testCred.Users = nil
	suite.adapter.On("CreateAccount", "test_login", suite.testCred).Return(nil)
	suite.adapter.On("SaveToken", mock.Anything).Return(nil)

	token, err := suite.service.Register("test_login", "test_password")
	suite.Require().NoError(err)
	suite.Equal("test_username", token.User.Username)
	suite.NotEmpty(token.Token)
	suite.Len(token.Token, 44)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestRegisterInvalidInput() {
	_, err := suite.service.Register("", "test_password")
	suite.Require().ErrorIs(err, types.ErrLoginEmpty)

	_, err = suite.service.Register("test_login", "")
	suite.Require().ErrorIs(err, types.ErrPasswordEmpty)

	_, err = suite.service.Register("test_login", "short")
	suite.Require().ErrorIs(err, types.ErrPasswordLength)
}

func (suite *AuthServiceTestSuite) TestLogout() {
	suite.adapter.On("RemoveToken", "valid_token").Return(nil)

	err := suite.service.Logout("valid_token")
	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestLogoutInvalidToken() {
	suite.adapter.On("RemoveToken", "invalid_token").Return(types.ErrTokenInvalid)

	err := suite.service.Logout("invalid_token")
	suite.Require().ErrorIs(err, types.ErrTokenInvalid)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestVerify() {
	suite.adapter.On("GetToken", "valid_token").Return(&types.Token{User: types.User{Username: "test_login"}}, nil)

	_, err := suite.service.Verify("valid_token")
	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestVerifyMasterKey() {
	suite.Require().Len(config.Current.MasterApiKey, 44)
	_, err := suite.service.Verify(config.Current.MasterApiKey)
	suite.Require().NoError(err)
}

func (suite *AuthServiceTestSuite) TestVerifyInvalidToken() {
	suite.adapter.On("GetToken", "invalid_token").Return(&types.Token{}, types.ErrTokenInvalid)

	_, err := suite.service.Verify("invalid_token")
	suite.Require().ErrorIs(err, types.ErrTokenInvalid)
	suite.adapter.AssertExpectations(suite.T())
}

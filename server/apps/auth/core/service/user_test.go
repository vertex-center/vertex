package service

import (
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/auth/core/port"
	"github.com/vertex-center/vertex/server/apps/auth/core/types"
)

type UserServiceTestSuite struct {
	suite.Suite

	service port.UserService
	adapter port.MockAuthAdapter

	testUser               types.User
	testCredentialsMethods []types.CredentialsMethods
	testErr                error
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupSubTest() {
	suite.adapter = port.MockAuthAdapter{}
	suite.service = NewUserService(&suite.adapter)

	suite.testUser = types.User{
		ID:       uuid.New(),
		Username: "test_username",
	}
	suite.testCredentialsMethods = []types.CredentialsMethods{{
		Name: types.CredentialsTypeLoginPassword,
	}}
	suite.testErr = errors.New("internal error")
}

func (suite *UserServiceTestSuite) TestGetUser() {
	suite.Run("success", func() {
		suite.adapter.On("GetUser", suite.testUser.Username).Return(suite.testUser, nil)
		user, err := suite.service.GetUser(suite.testUser.Username)
		suite.Require().NoError(err)
		suite.Equal(user, suite.testUser)
	})

	suite.Run("error", func() {
		suite.adapter.On("GetUser", suite.testUser.Username).Return(types.User{}, suite.testErr)
		user, err := suite.service.GetUser(suite.testUser.Username)
		suite.Require().ErrorIs(err, suite.testErr)
		suite.Empty(user)
	})
}

func (suite *UserServiceTestSuite) TestGetUserByID() {
	suite.Run("success", func() {
		suite.adapter.On("GetUserByID", suite.testUser.ID).Return(suite.testUser, nil)
		user, err := suite.service.GetUserByID(suite.testUser.ID)
		suite.Require().NoError(err)
		suite.Equal(user, suite.testUser)
	})

	suite.Run("error", func() {
		suite.adapter.On("GetUserByID", suite.testUser.ID).Return(types.User{}, suite.testErr)
		user, err := suite.service.GetUserByID(suite.testUser.ID)
		suite.Require().ErrorIs(err, suite.testErr)
		suite.Empty(user)
	})
}

func (suite *UserServiceTestSuite) TestPatchUser() {
	suite.Run("success", func() {
		newUser := suite.testUser
		newUser.Username = "new_username"
		suite.adapter.On("PatchUser", suite.testUser).Return(newUser, nil)
		user, err := suite.service.PatchUser(suite.testUser)
		suite.Require().NoError(err)
		suite.Equal(user, newUser)
	})

	suite.Run("error", func() {
		newUser := suite.testUser
		newUser.Username = "new_username"
		suite.adapter.On("PatchUser", suite.testUser).Return(types.User{}, suite.testErr)
		user, err := suite.service.PatchUser(suite.testUser)
		suite.Require().ErrorIs(err, suite.testErr)
		suite.Empty(user)
	})
}

func (suite *UserServiceTestSuite) TestGetUserCredentialsMethods() {
	suite.Run("success", func() {
		suite.adapter.On("GetUserCredentialsMethods", suite.testUser.ID).Return(suite.testCredentialsMethods, nil)
		creds, err := suite.service.GetUserCredentialsMethods(suite.testUser.ID)
		suite.Require().NoError(err)
		suite.Equal(creds, suite.testCredentialsMethods)
	})

	suite.Run("error", func() {
		suite.adapter.On("GetUserCredentialsMethods", suite.testUser.ID).Return([]types.CredentialsMethods{}, suite.testErr)
		creds, err := suite.service.GetUserCredentialsMethods(suite.testUser.ID)
		suite.Require().ErrorIs(err, suite.testErr)
		suite.Empty(creds)
	})
}

package service

import (
	"testing"
	"time"

	"github.com/juju/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type EmailTestSuite struct {
	suite.Suite

	testEmail types.Email
	adapter   port.MockEmailAdapter
	service   port.EmailService
}

func TestEmailTestSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

func (suite *EmailTestSuite) SetupSubTest() {
	suite.testEmail = types.Email{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Email:     "test@example.com",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	suite.adapter = port.MockEmailAdapter{}
	suite.service = NewEmailService(&suite.adapter)
}

func (suite *EmailTestSuite) TestGetEmails() {
	suite.Run("ok", func() {
		suite.adapter.On("GetEmails", suite.testEmail.UserID).Return([]types.Email{suite.testEmail}, nil)
		emails, err := suite.service.GetEmails(suite.testEmail.UserID)
		suite.Require().NoError(err)
		suite.Equal([]types.Email{suite.testEmail}, emails)
	})

	suite.Run("error", func() {
		errNotFound := errors.NotFoundf("email")
		suite.adapter.On("GetEmails", suite.testEmail.UserID).Return([]types.Email{}, errNotFound)
		emails, err := suite.service.GetEmails(suite.testEmail.UserID)
		suite.Require().ErrorIs(err, errNotFound)
		suite.Equal([]types.Email{}, emails)
	})
}

func (suite *EmailTestSuite) TestCreateEmail() {
	suite.Run("ok", func() {
		suite.adapter.On("CreateEmail", mock.AnythingOfType("*types.Email")).Return(nil)
		email, err := suite.service.CreateEmail(suite.testEmail.UserID, suite.testEmail.Email)
		suite.Require().NoError(err)
		suite.Equal(suite.testEmail.Email, email.Email)
		suite.Equal(suite.testEmail.UserID, email.UserID)
	})

	suite.Run("invalid email", func() {
		email, err := suite.service.CreateEmail(suite.testEmail.UserID, "invalid")
		suite.Require().Error(err)
		suite.Equal("create email address: mail: missing '@' or angle-addr", err.Error())
		suite.Equal(types.Email{}, email)
	})

	suite.Run("error", func() {
		err := errors.AlreadyExistsf("email")
		suite.adapter.On("CreateEmail", mock.AnythingOfType("*types.Email")).Return(err)
		email, err := suite.service.CreateEmail(suite.testEmail.UserID, suite.testEmail.Email)
		suite.Require().ErrorIs(err, err)
		suite.Equal(types.Email{}, email)
	})
}

func (suite *EmailTestSuite) TestDeleteEmail() {
	suite.Run("ok", func() {
		suite.adapter.On("DeleteEmail", suite.testEmail.UserID, suite.testEmail.Email).Return(nil)
		err := suite.service.DeleteEmail(suite.testEmail.UserID, suite.testEmail.Email)
		suite.Require().NoError(err)
	})

	suite.Run("error", func() {
		err := errors.NotFoundf("email")
		suite.adapter.On("DeleteEmail", suite.testEmail.UserID, suite.testEmail.Email).Return(err)
		err = suite.service.DeleteEmail(suite.testEmail.UserID, suite.testEmail.Email)
		suite.Require().ErrorIs(err, err)
	})
}

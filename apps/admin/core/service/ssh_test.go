package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/admin/core/port"
)

type SshServiceTestSuite struct {
	suite.Suite

	service *sshService
	adapter *port.MockSshAdapter
}

func TestSshServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SshServiceTestSuite))
}

func (suite *SshServiceTestSuite) SetupTest() {
	suite.adapter = &port.MockSshAdapter{}
	suite.service = NewSshService(suite.adapter).(*sshService)
}

func (suite *SshServiceTestSuite) TestGetAll() {
	suite.adapter.On("GetAll", context.Background()).Return(testDataAuthorizedKeys, nil)

	keys, err := suite.service.GetAll(context.Background())

	suite.Require().NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshServiceTestSuite) TestAdd() {
	suite.adapter.On("Add", context.Background(), testDataAuthorizedKey, "username").Return(nil)

	err := suite.service.Add(context.Background(), testDataAuthorizedKey, "username")

	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshServiceTestSuite) TestDelete() {
	suite.adapter.On("Remove", context.Background(), testDataFingerprint, "username").Return(nil)

	err := suite.service.Delete(context.Background(), testDataFingerprint, "username")

	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

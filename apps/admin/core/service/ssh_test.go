package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/admin/core/port"
)

type SshServiceTestSuite struct {
	suite.Suite

	service *SshService
	adapter *port.MockSshAdapter
}

func TestSshServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SshServiceTestSuite))
}

func (suite *SshServiceTestSuite) SetupTest() {
	suite.adapter = &port.MockSshAdapter{}
	suite.service = NewSshService(suite.adapter).(*SshService)
}

func (suite *SshServiceTestSuite) TestGetAll() {
	suite.adapter.On("GetAll").Return(testDataAuthorizedKeys, nil)

	keys, err := suite.service.GetAll()

	suite.Require().NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshServiceTestSuite) TestAdd() {
	suite.adapter.On("Add", testDataAuthorizedKey, "username").Return(nil)

	err := suite.service.Add(testDataAuthorizedKey, "username")

	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshServiceTestSuite) TestDelete() {
	suite.adapter.On("Remove", testDataFingerprint, "username").Return(nil)

	err := suite.service.Delete(testDataFingerprint, "username")

	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SshServiceTestSuite struct {
	suite.Suite

	service SshService
	adapter MockSshAdapter
}

func TestSshServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SshServiceTestSuite))
}

func (suite *SshServiceTestSuite) SetupSuite() {
	suite.adapter = MockSshAdapter{}
	suite.service = NewSshService(&suite.adapter)
}

func (suite *SshServiceTestSuite) TestGetAll() {
	suite.adapter.On("GetAll").Return(testDataAuthorizedKeys, nil)

	keys, err := suite.service.GetAll()

	suite.NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshServiceTestSuite) TestAdd() {
	suite.adapter.On("Add", testDataAuthorizedKey).Return(nil)

	err := suite.service.Add(testDataAuthorizedKey)

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshServiceTestSuite) TestDelete() {
	suite.adapter.On("Remove", testDataFingerprint).Return(nil)

	err := suite.service.Delete(testDataFingerprint)

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

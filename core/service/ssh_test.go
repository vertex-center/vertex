package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
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
	suite.adapter.GetAllFunc = func() ([]types.PublicKey, error) {
		return testDataAuthorizedKeys, nil
	}

	keys, err := suite.service.GetAll()

	suite.Require().NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.Equal(1, suite.adapter.GetAllCalls)
}

func (suite *SshServiceTestSuite) TestAdd() {
	suite.adapter.AddFunc = func(key string) error {
		return nil
	}

	err := suite.service.Add(testDataAuthorizedKey)

	suite.Require().NoError(err)
	suite.Equal(1, suite.adapter.AddCalls)
}

func (suite *SshServiceTestSuite) TestDelete() {
	suite.adapter.RemoveFunc = func(fingerprint string) error {
		return nil
	}

	err := suite.service.Delete(testDataFingerprint)

	suite.Require().NoError(err)
	suite.Equal(1, suite.adapter.RemoveCalls)
}

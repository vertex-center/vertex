package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

const (
	testDataAuthorizedKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6IPH4bqdhPVQUfdmuisPdQJO6Tv2+a0OZ9qLs6W0W2flxn6/yQmYut02cl0UtNcDtmb4RqNj2ms2v2TeDVSWVZkUR/q4jjZSSljQEpTd3r1YhYrO/GPDNiIUMm5HvZ8qIfBQA6gn9uMT1g6FO53O64ACNr+ItU4gNdr+S44MNJRMxMy6+s/LsFlQjyO2MbPQHQ6HSOgTLrCNiH8NTLA/evekrZ/rmIZrrES2vQvw5pbCDgEOkLZruRSMMFJFStb6tlGoiN/jQpfX51jebDVLZ1/U3SU5+7LNN6DxZYE9w1eCA2G8L8q1PUYju+b4F6IhGA1AYXPaAaR12qRJ4lLeN"
	testDataFingerprint   = "SHA256:eLfsDB1H1SrvT7Bgo9U1i/ATcldIrOqin2H0MGEy5I8"
)

var (
	testDataAuthorizedKeys = []types.PublicKey{
		{
			Type:              "ssh-rsa",
			FingerprintSHA256: "SHA256:eLfsDB1H1SrvT7Bgo9U1i/ATcldIrOqin2H0MGEy5I8",
			Username:          "username",
		},
		{
			Type:              "ssh-rsa",
			FingerprintSHA256: "SHA256:ubvRPPaAlkFeuFQeC748c43nRPTjaRGxnG9C0j+WlJ0",
			Username:          "username",
		},
	}
)

type SshKernelServiceTestSuite struct {
	suite.Suite

	testUser  user.User
	testUsers []user.User

	service *SshKernelService
	adapter *port.MockSshKernelAdapter
}

func TestSshKernelServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SshKernelServiceTestSuite))
}

func (suite *SshKernelServiceTestSuite) SetupTest() {
	suite.testUser = user.User{
		Name: "username",
	}
	suite.testUsers = []user.User{suite.testUser}
	suite.adapter = &port.MockSshKernelAdapter{}
	suite.service = NewSshKernelService(suite.adapter).(*SshKernelService)
}

func (suite *SshKernelServiceTestSuite) TestGetAll() {
	suite.adapter.On("GetAll", suite.testUsers).Return(testDataAuthorizedKeys, nil)
	suite.adapter.On("GetUsers").Return(suite.testUsers, nil)

	keys, err := suite.service.GetAll()

	suite.Require().NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshKernelServiceTestSuite) TestAdd() {
	suite.adapter.On("Add", testDataAuthorizedKey, suite.testUser).Return(nil)
	suite.adapter.On("GetUsers").Return(suite.testUsers, nil)

	err := suite.service.Add(testDataAuthorizedKey, suite.testUser.Name)

	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshKernelServiceTestSuite) TestAddInvalidKey() {
	suite.adapter.On("Add", "invalid", suite.testUser).Return(nil)
	suite.adapter.On("GetUsers").Return(suite.testUsers, nil)

	err := suite.service.Add("invalid", "username")

	suite.Require().Error(err)
	suite.Require().ErrorIsf(err, ErrInvalidPublicKey, "invalid key")
	suite.adapter.AssertNotCalled(suite.T(), "Add", "invalid")
}

func (suite *SshKernelServiceTestSuite) TestDelete() {
	suite.adapter.On("Remove", testDataFingerprint, suite.testUser).Return(nil)
	suite.adapter.On("GetUsers").Return(suite.testUsers, nil)

	err := suite.service.Delete(testDataFingerprint, "username")

	suite.Require().NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

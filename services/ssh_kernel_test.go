package services

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
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
		},
		{
			Type:              "ssh-rsa",
			FingerprintSHA256: "SHA256:ubvRPPaAlkFeuFQeC748c43nRPTjaRGxnG9C0j+WlJ0",
		},
	}
)

type SshKernelServiceTestSuite struct {
	suite.Suite

	service SshKernelService
	adapter MockSshAdapter
}

func TestSshKernelServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SshKernelServiceTestSuite))
}

func (suite *SshKernelServiceTestSuite) SetupSuite() {
	suite.adapter = MockSshAdapter{}
	suite.service = NewSshKernelService(&suite.adapter)
}

func (suite *SshKernelServiceTestSuite) TestGetAll() {
	suite.adapter.On("GetAll").Return(testDataAuthorizedKeys, nil)

	keys, err := suite.service.GetAll()

	suite.NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshKernelServiceTestSuite) TestAdd() {
	suite.adapter.On("Add", testDataAuthorizedKey).Return(nil)

	err := suite.service.Add(testDataAuthorizedKey)

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

func (suite *SshKernelServiceTestSuite) TestAddInvalidKey() {
	suite.adapter.On("Add", testDataAuthorizedKey).Return(nil)

	err := suite.service.Add("invalid")

	suite.Error(err)
	suite.ErrorIsf(err, ErrInvalidPublicKey, "invalid key")
	suite.adapter.AssertNotCalled(suite.T(), "Add", "invalid")
}

func (suite *SshKernelServiceTestSuite) TestDelete() {
	suite.adapter.On("Remove", testDataFingerprint).Return(nil)

	err := suite.service.Delete(testDataFingerprint)

	suite.NoError(err)
	suite.adapter.AssertExpectations(suite.T())
}

type MockSshAdapter struct {
	mock.Mock
}

func (m *MockSshAdapter) GetAll() ([]types.PublicKey, error) {
	args := m.Called()
	return args.Get(0).([]types.PublicKey), args.Error(1)
}

func (m *MockSshAdapter) Add(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockSshAdapter) Remove(fingerprint string) error {
	args := m.Called(fingerprint)
	return args.Error(0)
}

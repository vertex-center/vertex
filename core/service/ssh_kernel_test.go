package service

import (
	"testing"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"

	"github.com/stretchr/testify/suite"
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

	service *SshKernelService
	adapter *port.MockSshAdapter
}

func TestSshKernelServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SshKernelServiceTestSuite))
}

func (suite *SshKernelServiceTestSuite) SetupTest() {
	suite.adapter = &port.MockSshAdapter{}
	suite.service = NewSshKernelService(suite.adapter).(*SshKernelService)
}

func (suite *SshKernelServiceTestSuite) TestGetAll() {
	suite.adapter.GetAllFunc = func() ([]types.PublicKey, error) {
		return testDataAuthorizedKeys, nil
	}

	keys, err := suite.service.GetAll()

	suite.Require().NoError(err)
	suite.Equal(testDataAuthorizedKeys, keys)
	suite.Equal(1, suite.adapter.GetAllCalls)
}

func (suite *SshKernelServiceTestSuite) TestAdd() {
	suite.adapter.AddFunc = func(key string) error {
		return nil
	}

	err := suite.service.Add(testDataAuthorizedKey)

	suite.Require().NoError(err)
	suite.Equal(1, suite.adapter.AddCalls)
}

func (suite *SshKernelServiceTestSuite) TestAddInvalidKey() {
	suite.adapter.AddFunc = func(key string) error {
		return nil
	}

	err := suite.service.Add("invalid")

	suite.Require().Error(err)
	suite.Require().ErrorIsf(err, ErrInvalidPublicKey, "invalid key")
	suite.Equal(0, suite.adapter.AddCalls)
}

func (suite *SshKernelServiceTestSuite) TestDelete() {
	suite.adapter.RemoveFunc = func(fingerprint string) error {
		return nil
	}

	err := suite.service.Delete(testDataFingerprint)

	suite.Require().NoError(err)
	suite.Equal(1, suite.adapter.RemoveCalls)
}

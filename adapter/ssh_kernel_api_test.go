package adapter

import (
	"net/http"
	"testing"

	"github.com/vertex-center/vertex/core/types"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/config"
)

type SshKernelApiAdapterTestSuite struct {
	suite.Suite

	adapter SshKernelApiAdapter
}

func TestSshKernelApiAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(SshKernelApiAdapterTestSuite))
}

func (suite *SshKernelApiAdapterTestSuite) SetupTest() {
	suite.adapter = *NewSshKernelApiAdapter().(*SshKernelApiAdapter)
}

func (suite *SshKernelApiAdapterTestSuite) TestGetAll() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		Get("/api/security/ssh").
		Reply(http.StatusOK).
		JSON([]types.PublicKey{})

	keys, err := suite.adapter.GetAll()
	suite.Require().NoError(err)
	suite.Empty(keys)
}

func (suite *SshKernelApiAdapterTestSuite) TestAdd() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		Post("/api/security/ssh").
		Reply(http.StatusOK)

	err := suite.adapter.Add("key")
	suite.Require().NoError(err)
}

func (suite *SshKernelApiAdapterTestSuite) TestDelete() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		Delete("/api/security/ssh/fingerprint").
		Reply(http.StatusOK)

	err := suite.adapter.Remove("fingerprint")
	suite.Require().NoError(err)
}

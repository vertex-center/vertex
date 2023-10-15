package adapter

import (
	"github.com/vertex-center/vertex/core/types"
	"net/http"
	"testing"

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
	suite.NoError(err)
	suite.Len(keys, 0)
}

func (suite *SshKernelApiAdapterTestSuite) TestAdd() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		Post("/api/security/ssh").
		Reply(http.StatusOK)

	err := suite.adapter.Add("key")
	suite.NoError(err)
}

func (suite *SshKernelApiAdapterTestSuite) TestDelete() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		Delete("/api/security/ssh/fingerprint").
		Reply(http.StatusOK)

	err := suite.adapter.Remove("fingerprint")
	suite.NoError(err)
}

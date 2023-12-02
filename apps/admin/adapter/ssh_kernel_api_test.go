package adapter

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/user"
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
		Get("/api/app/admin/ssh").
		Reply(http.StatusOK).
		JSON([]types.PublicKey{})

	keys, err := suite.adapter.GetAll()
	suite.Require().NoError(err)
	suite.Empty(keys)
}

func (suite *SshKernelApiAdapterTestSuite) TestAdd() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		JSON(&handler.AddSSHKeyBody{
			AuthorizedKey: "key",
			Username:      "username",
		}).
		Post("/api/app/admin/ssh").
		Reply(http.StatusOK)

	err := suite.adapter.Add("key", "username")
	suite.Require().NoError(err)
}

func (suite *SshKernelApiAdapterTestSuite) TestDelete() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		JSON(&handler.DeleteSSHKeyBody{
			Fingerprint: "fingerprint",
			Username:    "username",
		}).
		Delete("/api/app/admin/ssh").
		Reply(http.StatusOK)

	err := suite.adapter.Remove("fingerprint", "username")
	suite.Require().NoError(err)
}

func (suite *SshKernelApiAdapterTestSuite) TestGetUsers() {
	gock.Off()
	gock.New(config.Current.KernelURL()).
		Get("/api/app/admin/ssh/users").
		Reply(http.StatusOK).
		JSON([]user.User{})

	users, err := suite.adapter.GetUsers()
	suite.Require().NoError(err)
	suite.Empty(users)
}

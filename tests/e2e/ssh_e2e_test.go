//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/api"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/user"
)

type SshE2ETestSuite struct {
	suite.Suite
}

func TestE2ESshTestSuite(t *testing.T) {
	suite.Run(t, new(SshE2ETestSuite))
}

func (suite *SshE2ETestSuite) SetupSuite() {
	config.Current.Port = "7130"
	config.Current.PortKernel = "7131"
}

func (suite *SshE2ETestSuite) TestSsh() {
	sshClient := api.NewClient()

	ctx := context.Background()

	users, err := sshClient.GetSSHUsers(ctx)
	suite.Require().NoError(err)
	suite.Len(users, 2)
	suite.Equal([]user.User{
		{
			Name:    "root",
			HomeDir: "/",
		},
		{
			Name:    "nexa",
			HomeDir: "/home/nexa",
		},
	}, users)
}

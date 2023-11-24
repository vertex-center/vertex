//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/api"
	"github.com/vertex-center/vertex/config"
	"golang.org/x/crypto/ssh"
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

	// Get SSH users
	users, err := sshClient.GetSSHUsers(ctx)
	suite.Require().NoError(err)
	suite.Equal([]string{"nexa"}, users)

	// Get SSH keys
	keys, err := sshClient.GetSSHKeys(ctx)
	suite.Require().NoError(err)
	suite.Empty(keys)

	// Add SSH key
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6IPH4bqdhPVQUfdmuisPdQJO6Tv2+a0OZ9qLs6W0W2flxn6/yQmYut02cl0UtNcDtmb4RqNj2ms2v2TeDVSWVZkUR/q4jjZSSljQEpTd3r1YhYrO/GPDNiIUMm5HvZ8qIfBQA6gn9uMT1g6FO53O64ACNr+ItU4gNdr+S44MNJRMxMy6+s/LsFlQjyO2MbPQHQ6HSOgTLrCNiH8NTLA/evekrZ/rmIZrrES2vQvw5pbCDgEOkLZruRSMMFJFStb6tlGoiN/jQpfX51jebDVLZ1/U3SU5+7LNN6DxZYE9w1eCA2G8L8q1PUYju+b4F6IhGA1AYXPaAaR12qRJ4lLeN"
	err = sshClient.AddSSHKey(ctx, key, "nexa")
	suite.Require().NoError(err)

	// Check if key was added
	keys, err = sshClient.GetSSHKeys(ctx)
	suite.Require().NoError(err)
	suite.Len(keys, 1)

	// Delete SSH key
	k, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(key))
	fingerprint := ssh.FingerprintSHA256(k)
	err = sshClient.DeleteSSHKey(ctx, fingerprint, "nexa")
	suite.Require().NoError(err)

	// Check if key was deleted
	keys, err = sshClient.GetSSHKeys(ctx)
	suite.Require().NoError(err)
	suite.Empty(keys)
}

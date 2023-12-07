package adapter

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/pkg/user"
	"golang.org/x/crypto/ssh"
)

var (
	keys = []string{
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6IPH4bqdhPVQUfdmuisPdQJO6Tv2+a0OZ9qLs6W0W2flxn6/yQmYut02cl0UtNcDtmb4RqNj2ms2v2TeDVSWVZkUR/q4jjZSSljQEpTd3r1YhYrO/GPDNiIUMm5HvZ8qIfBQA6gn9uMT1g6FO53O64ACNr+ItU4gNdr+S44MNJRMxMy6+s/LsFlQjyO2MbPQHQ6HSOgTLrCNiH8NTLA/evekrZ/rmIZrrES2vQvw5pbCDgEOkLZruRSMMFJFStb6tlGoiN/jQpfX51jebDVLZ1/U3SU5+7LNN6DxZYE9w1eCA2G8L8q1PUYju+b4F6IhGA1AYXPaAaR12qRJ4lLeN",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCtkVmRevgiIRc7QHahcd01d+0qjtZj1KcY5u25TQW7GomgVuJukdKupnUP2Q1DGo1JjI0OMaIVcEAs4rQgHDAIYovHSeQpkhb3QzJKpS9YUxq/ZWtBQd7cdyRrwAJuT0uR0m52NopEVaaETSIFH6byScRoOAdKgRPwWv5EiHleklOuZCG2/BKq2FtHIb5xb7eAEeMy/5ebu1f4C211/q/Y0AIy/Gp7rJGTDSutTi2UXMQxo3kVDykIIg/xqH2h5IUvYOR8Y+t6f9rbKPcglc+9ygmYHeqrIVmkFzru1sbOOCHlIfv1N53RVp5A9734cHm9u3FzfIPkV+j0tOJ8dhdP",
	}

	fingerprints = []string{
		"SHA256:eLfsDB1H1SrvT7Bgo9U1i/ATcldIrOqin2H0MGEy5I8",
		"SHA256:ubvRPPaAlkFeuFQeC748c43nRPTjaRGxnG9C0j+WlJ0",
	}
)

type SshFsAdapterTestSuite struct {
	suite.Suite

	testUser user.User

	adapter            sshFsAdapter
	authorizedKeysFile *os.File
}

func TestSshFsAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(SshFsAdapterTestSuite))
}

func (suite *SshFsAdapterTestSuite) SetupTest() {
	var err error

	suite.testUser = user.User{
		Name:    "test",
		HomeDir: suite.T().TempDir(),
	}

	err = os.MkdirAll(path.Join(suite.testUser.HomeDir, ".ssh"), os.ModePerm)
	if err != nil {
		suite.FailNow(err.Error())
	}

	suite.authorizedKeysFile, err = os.Create(path.Join(suite.testUser.HomeDir, ".ssh", "authorized_keys"))
	if err != nil {
		suite.FailNow(err.Error())
	}

	for i := range keys {
		_, err = suite.authorizedKeysFile.WriteString(keys[i] + "\n")
		if err != nil {
			suite.FailNow(err.Error())
		}
	}

	suite.adapter = *NewSshFsAdapter().(*sshFsAdapter)
}

func (suite *SshFsAdapterTestSuite) TearDownTest() {
	suite.authorizedKeysFile.Close()

	err := os.Remove(suite.authorizedKeysFile.Name())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		suite.Require().NoError(err)
	}

	err = os.RemoveAll(suite.testUser.HomeDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		suite.Require().NoError(err)
	}
}

func (suite *SshFsAdapterTestSuite) TestGetAll() {
	keys, err := suite.adapter.GetAll([]user.User{suite.testUser})
	suite.Require().NoError(err)
	suite.Len(keys, 2)
	for i, key := range keys {
		suite.Equal("ssh-rsa", key.Type)
		suite.Equal(fingerprints[i], key.FingerprintSHA256)
	}
}

func (suite *SshFsAdapterTestSuite) TestGetAllInvalidKey() {
	_, err := suite.authorizedKeysFile.Write([]byte("invalid"))
	suite.Require().NoError(err)

	keys, err := suite.adapter.GetAll([]user.User{suite.testUser})
	suite.Require().NoError(err)
	suite.Len(keys, 2)
}

func (suite *SshFsAdapterTestSuite) TestGetAllNoSuchFile() {
	suite.authorizedKeysFile.Close()
	err := os.Remove(path.Join(suite.testUser.HomeDir, ".ssh", "authorized_keys"))
	suite.Require().NoError(err)

	keys, err := suite.adapter.GetAll([]user.User{suite.testUser})
	suite.Require().NoError(err)
	suite.Empty(keys)
}

func (suite *SshFsAdapterTestSuite) TestAdd() {
	publicKey, err := generatePublicKey()
	if err != nil {
		suite.FailNow(err.Error())
	}

	err = suite.adapter.Add(string(publicKey), suite.testUser)
	suite.Require().NoError(err)

	keys, err := suite.adapter.GetAll([]user.User{suite.testUser})
	suite.Require().NoError(err)
	suite.Len(keys, 3)
}

func (suite *SshFsAdapterTestSuite) TestDelete() {
	k, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(keys[0]))
	err := suite.adapter.Remove(ssh.FingerprintSHA256(k), suite.testUser)
	suite.Require().NoError(err)

	keys, err := suite.adapter.GetAll([]user.User{suite.testUser})
	suite.Require().NoError(err)
	suite.Len(keys, 1)
}

func generatePublicKey() ([]byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	err = key.Validate()
	if err != nil {
		return nil, err
	}

	publicKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(publicKey), nil
}

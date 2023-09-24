package services

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/ssh"
)

type SSHServiceTestSuite struct {
	suite.Suite

	service            SSHService
	authorizedKeysFile *os.File

	// Test data
	keys           []string
	authorizedKeys []ssh.PublicKey
	fingerprints   []string
}

func TestSSHServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SSHServiceTestSuite))
}

func (suite *SSHServiceTestSuite) SetupSuite() {
	suite.keys = []string{
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6IPH4bqdhPVQUfdmuisPdQJO6Tv2+a0OZ9qLs6W0W2flxn6/yQmYut02cl0UtNcDtmb4RqNj2ms2v2TeDVSWVZkUR/q4jjZSSljQEpTd3r1YhYrO/GPDNiIUMm5HvZ8qIfBQA6gn9uMT1g6FO53O64ACNr+ItU4gNdr+S44MNJRMxMy6+s/LsFlQjyO2MbPQHQ6HSOgTLrCNiH8NTLA/evekrZ/rmIZrrES2vQvw5pbCDgEOkLZruRSMMFJFStb6tlGoiN/jQpfX51jebDVLZ1/U3SU5+7LNN6DxZYE9w1eCA2G8L8q1PUYju+b4F6IhGA1AYXPaAaR12qRJ4lLeN",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCtkVmRevgiIRc7QHahcd01d+0qjtZj1KcY5u25TQW7GomgVuJukdKupnUP2Q1DGo1JjI0OMaIVcEAs4rQgHDAIYovHSeQpkhb3QzJKpS9YUxq/ZWtBQd7cdyRrwAJuT0uR0m52NopEVaaETSIFH6byScRoOAdKgRPwWv5EiHleklOuZCG2/BKq2FtHIb5xb7eAEeMy/5ebu1f4C211/q/Y0AIy/Gp7rJGTDSutTi2UXMQxo3kVDykIIg/xqH2h5IUvYOR8Y+t6f9rbKPcglc+9ygmYHeqrIVmkFzru1sbOOCHlIfv1N53RVp5A9734cHm9u3FzfIPkV+j0tOJ8dhdP",
	}

	for _, key := range suite.keys {
		k, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(key))
		if k == nil {
			continue
		}
		suite.authorizedKeys = append(suite.authorizedKeys, k)
		suite.fingerprints = append(suite.fingerprints, ssh.FingerprintSHA256(k))
	}
}

func (suite *SSHServiceTestSuite) SetupTest() {
	var err error

	suite.authorizedKeysFile, err = os.CreateTemp("", "*_authorized_keys")
	if err != nil {
		suite.FailNow(err.Error())
	}

	for i := range suite.keys {
		_, err = suite.authorizedKeysFile.WriteString(suite.keys[i] + "\n")
		if err != nil {
			suite.FailNow(err.Error())
		}
	}

	suite.service = NewSSHService(&SSHServiceParams{
		AuthorizedKeysPath: suite.authorizedKeysFile.Name(),
	})
}

func (suite *SSHServiceTestSuite) TearDownTest() {
	err := os.Remove(suite.authorizedKeysFile.Name())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		suite.NoError(err)
	}
}

func (suite *SSHServiceTestSuite) TestGetAll() {
	keys, err := suite.service.GetAll()
	suite.NoError(err)
	suite.Equal(2, len(keys))
	for i, key := range keys {
		suite.Equal("ssh-rsa", key.Type)
		suite.Equal(suite.fingerprints[i], key.FingerprintSHA256)
	}
}

func (suite *SSHServiceTestSuite) TestGetAllInvalidKey() {
	_, err := suite.authorizedKeysFile.Write([]byte("invalid"))
	suite.NoError(err)

	keys, err := suite.service.GetAll()
	suite.NoError(err)
	suite.Equal(2, len(keys))
}

func (suite *SSHServiceTestSuite) TestGetAllNoSuchFile() {
	suite.authorizedKeysFile.Close()
	err := os.Remove(suite.service.authorizedKeysPath)
	suite.NoError(err)

	keys, err := suite.service.GetAll()
	suite.NoError(err)
	suite.Equal(0, len(keys))
}

func (suite *SSHServiceTestSuite) TestAdd() {
	publicKey, err := generatePublicKey()
	if err != nil {
		suite.FailNow(err.Error())
	}

	err = suite.service.Add(string(publicKey))
	suite.NoError(err)

	keys, err := suite.service.GetAll()
	suite.NoError(err)
	suite.Equal(3, len(keys))
}

func (suite *SSHServiceTestSuite) TestAddInvalidKey() {
	err := suite.service.Add("invalid")
	suite.Error(err)

	keys, err := suite.service.GetAll()
	suite.NoError(err)
	suite.Equal(2, len(keys))
}

func (suite *SSHServiceTestSuite) TestDelete() {
	err := suite.service.Delete(ssh.FingerprintSHA256(suite.authorizedKeys[0]))
	suite.NoError(err)

	keys, err := suite.service.GetAll()
	suite.NoError(err)
	suite.Equal(1, len(keys))
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

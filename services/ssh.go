package services

import (
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"golang.org/x/crypto/ssh"
)

type SSHService struct{}

func NewSSHService() SSHService {
	return SSHService{}
}

func (s *SSHService) GetAll() ([]types.PublicKey, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	authorizedKeysPath := path.Join(dir, ".ssh", "authorized_keys")

	bytes, err := os.ReadFile(authorizedKeysPath)
	if err != nil {
		return nil, err
	}

	var publicKeys []ssh.PublicKey
	for len(bytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(bytes)
		if err != nil {
			log.Error(err)
		}
		publicKeys = append(publicKeys, pubKey)
		bytes = rest
	}

	var keys []types.PublicKey
	for _, key := range publicKeys {
		keys = append(keys, types.PublicKey{
			Type:              key.Type(),
			FingerprintSHA256: ssh.FingerprintSHA256(key),
		})
	}

	return keys, nil
}

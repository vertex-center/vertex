package services

import (
	"errors"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"golang.org/x/crypto/ssh"
)

var (
	ErrInvalidPublicKey = errors.New("invalid key")
)

type SSHService struct{}

func NewSSHService() SSHService {
	return SSHService{}
}

func (s *SSHService) GetAll() ([]types.PublicKey, error) {
	authorizedKeysPath, err := s.getAuthorizedKeysPath()
	if err != nil {
		return nil, err
	}

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

// Add adds an SSH key to the authorized keys file. The key must
// be a valid SSH public key, otherwise ErrInvalidPublicKey is returned.
func (s *SSHService) Add(authorizedKey string) error {
	// Check if the key is valid.
	_, err := ssh.ParsePublicKey([]byte(authorizedKey))
	if err != nil {
		return ErrInvalidPublicKey
	}

	authorizedKeysPath, err := s.getAuthorizedKeysPath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(authorizedKeysPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(authorizedKey + "\n")
	return err
}

func (s *SSHService) getAuthorizedKeysPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ".ssh", "authorized_keys"), nil
}

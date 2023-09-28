package services

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"golang.org/x/crypto/ssh"
)

var (
	ErrInvalidPublicKey = errors.New("invalid key")
)

type SSHKernelService struct {
	authorizedKeysPath string
}

type SSHKernelServiceParams struct {
	AuthorizedKeysPath string
}

func NewSSHKernelService(params *SSHKernelServiceParams) SSHKernelService {
	s := SSHKernelService{}

	if params == nil {
		params = &SSHKernelServiceParams{}
	}

	s.authorizedKeysPath = params.AuthorizedKeysPath
	if s.authorizedKeysPath == "" {
		var err error
		s.authorizedKeysPath, err = getAuthorizedKeysPath()
		if err != nil {
			log.Error(err)
		}
	}

	return s
}

func (s *SSHKernelService) GetAll() ([]types.PublicKey, error) {
	bytes, err := os.ReadFile(s.authorizedKeysPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Info("authorized_keys file does not exist")
		return []types.PublicKey{}, nil
	} else if err != nil {
		return nil, err
	}

	var publicKeys []ssh.PublicKey
	for len(bytes) > 0 {
		pubKey, _, _, rest, _ := ssh.ParseAuthorizedKey(bytes)
		if pubKey != nil {
			publicKeys = append(publicKeys, pubKey)
		}
		bytes = rest
	}

	keys := []types.PublicKey{}
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
func (s *SSHKernelService) Add(authorizedKey string) error {
	// Check if the key is valid.
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authorizedKey))
	if err != nil {
		return ErrInvalidPublicKey
	}

	file, err := os.OpenFile(s.authorizedKeysPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(authorizedKey + "\n")
	return err
}

// Delete deletes an SSH key from the authorized keys file.
func (s *SSHKernelService) Delete(fingerprint string) error {
	content, err := os.ReadFile(s.authorizedKeysPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(line))
		if key == nil {
			continue
		}

		fingerprintLine := ssh.FingerprintSHA256(key)

		if fingerprintLine == fingerprint {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	return os.WriteFile(s.authorizedKeysPath, []byte(strings.Join(lines, "\n")), 0644)
}

func getAuthorizedKeysPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ".ssh", "authorized_keys"), nil
}

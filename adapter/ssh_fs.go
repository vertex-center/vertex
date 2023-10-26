package adapter

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"

	"github.com/vertex-center/vertex/pkg/log"
	"golang.org/x/crypto/ssh"
)

type SshFsAdapter struct {
	authorizedKeysPath string
}

type SshFsAdapterParams struct {
	AuthorizedKeysPath string
}

func NewSshFsAdapter(params *SshFsAdapterParams) port.SshAdapter {
	s := &SshFsAdapter{}

	if params == nil {
		params = &SshFsAdapterParams{}
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

func (a *SshFsAdapter) GetAll() ([]types.PublicKey, error) {
	bytes, err := os.ReadFile(a.authorizedKeysPath)
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

func (a *SshFsAdapter) Add(key string) error {
	file, err := os.OpenFile(a.authorizedKeysPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(key + "\n")
	return err
}

func (a *SshFsAdapter) Remove(fingerprint string) error {
	content, err := os.ReadFile(a.authorizedKeysPath)
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

	return os.WriteFile(a.authorizedKeysPath, []byte(strings.Join(lines, "\n")), 0644)
}

func getAuthorizedKeysPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ".ssh", "authorized_keys"), nil
}

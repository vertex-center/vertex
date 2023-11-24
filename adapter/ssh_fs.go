package adapter

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/user"

	"github.com/vertex-center/vertex/pkg/log"
	"golang.org/x/crypto/ssh"
)

type SshFsAdapter struct {
}

func NewSshFsAdapter() port.SshKernelAdapter {
	return &SshFsAdapter{}
}

func (a *SshFsAdapter) GetAll(users []user.User) ([]types.PublicKey, error) {
	var keys []types.PublicKey
	for _, u := range users {
		p := getAuthorizedKeysPath(u.HomeDir)

		bytes, err := os.ReadFile(p)
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

		for _, key := range publicKeys {
			keys = append(keys, types.PublicKey{
				Type:              key.Type(),
				FingerprintSHA256: ssh.FingerprintSHA256(key),
				Username:          u.Name,
			})
		}
	}
	return keys, nil
}

func (a *SshFsAdapter) Add(key string, user user.User) error {
	p := getAuthorizedKeysPath(user.HomeDir)

	file, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(key + "\n")
	return err
}

func (a *SshFsAdapter) Remove(fingerprint string, user user.User) error {
	p := getAuthorizedKeysPath(user.HomeDir)

	content, err := os.ReadFile(p)
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

	return os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0644)
}

func (a *SshFsAdapter) GetUsers() ([]user.User, error) {
	return user.GetAll()
}

func getAuthorizedKeysPath(homeDir string) string {
	return path.Join(homeDir, ".ssh", "authorized_keys")
}

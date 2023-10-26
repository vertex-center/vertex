package service

import (
	"errors"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"golang.org/x/crypto/ssh"
)

var (
	ErrInvalidPublicKey = errors.New("invalid key")
)

type SshKernelService struct {
	sshAdapter port.SshAdapter
}

func NewSshKernelService(sshAdapter port.SshAdapter) port.SshService {
	return &SshKernelService{
		sshAdapter: sshAdapter,
	}
}

// GetAll returns all SSH keys from the authorized keys file.
func (s *SshKernelService) GetAll() ([]types.PublicKey, error) {
	return s.sshAdapter.GetAll()
}

// Add adds an SSH key to the authorized keys file. The key must
// be a valid SSH public key, otherwise ErrInvalidPublicKey is returned.
func (s *SshKernelService) Add(authorizedKey string) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authorizedKey))
	if err != nil {
		return ErrInvalidPublicKey
	}
	return s.sshAdapter.Add(authorizedKey)
}

// Delete deletes an SSH key from the authorized keys file.
func (s *SshKernelService) Delete(fingerprint string) error {
	return s.sshAdapter.Remove(fingerprint)
}

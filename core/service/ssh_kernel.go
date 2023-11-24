package service

import (
	"errors"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/user"
	"golang.org/x/crypto/ssh"
)

var (
	ErrInvalidPublicKey = errors.New("invalid key")
	ErrUserNotFound     = errors.New("user not found")
)

type SshKernelService struct {
	sshAdapter port.SshKernelAdapter
}

func NewSshKernelService(sshAdapter port.SshKernelAdapter) port.SshKernelService {
	return &SshKernelService{
		sshAdapter: sshAdapter,
	}
}

// GetAll returns all SSH keys from the authorized keys files.
func (s *SshKernelService) GetAll() ([]types.PublicKey, error) {
	users, err := s.sshAdapter.GetUsers()
	if err != nil {
		return nil, err
	}
	return s.sshAdapter.GetAll(users)
}

// Add adds an SSH key to the authorized keys file. The key must
// be a valid SSH public key, otherwise ErrInvalidPublicKey is returned.
func (s *SshKernelService) Add(authorizedKey string, username string) error {
	u, err := s.GetUser(username)
	if err != nil {
		return err
	}

	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(authorizedKey))
	if err != nil {
		return ErrInvalidPublicKey
	}
	return s.sshAdapter.Add(authorizedKey, u)
}

// Delete deletes an SSH key from the authorized keys file.
func (s *SshKernelService) Delete(fingerprint string, username string) error {
	u, err := s.GetUser(username)
	if err != nil {
		return err
	}

	return s.sshAdapter.Remove(fingerprint, u)
}

// GetUsers returns all users on the system that can have SSH keys.
func (s *SshKernelService) GetUsers() ([]user.User, error) {
	return s.sshAdapter.GetUsers()
}

func (s *SshKernelService) GetUser(username string) (user.User, error) {
	users, err := s.GetUsers()
	if err != nil {
		return user.User{}, err
	}

	for _, u := range users {
		if u.Name == username {
			return u, nil
		}
	}

	return user.User{}, ErrUserNotFound
}

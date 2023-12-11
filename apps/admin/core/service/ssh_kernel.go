package service

import (
	"context"
	"errors"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
	"golang.org/x/crypto/ssh"
)

var (
	ErrInvalidPublicKey = errors.New("invalid key")
	ErrUserNotFound     = errors.New("user not found")
	ErrFailedToAddKey   = errors.New("failed to add key")
)

type sshKernelService struct {
	sshAdapter port.SshKernelAdapter
}

func NewSshKernelService(sshAdapter port.SshKernelAdapter) port.SshKernelService {
	return &sshKernelService{
		sshAdapter: sshAdapter,
	}
}

// GetAll returns all SSH keys from the authorized keys files.
func (s *sshKernelService) GetAll() ([]types.PublicKey, error) {
	users, err := s.sshAdapter.GetUsers(context.Background())
	if err != nil {
		return nil, err
	}
	return s.sshAdapter.GetAll(context.Background(), users)
}

// Add adds an SSH key to the authorized keys file. The key must
// be a valid SSH public key, otherwise ErrInvalidPublicKey is returned.
func (s *sshKernelService) Add(authorizedKey string, username string) error {
	u, err := s.GetUser(username)
	if err != nil {
		return err
	}

	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(authorizedKey))
	if err != nil {
		return ErrInvalidPublicKey
	}
	return s.sshAdapter.Add(context.Background(), authorizedKey, u)
}

// Delete deletes an SSH key from the authorized keys file.
func (s *sshKernelService) Delete(fingerprint string, username string) error {
	u, err := s.GetUser(username)
	if err != nil {
		return err
	}

	return s.sshAdapter.Remove(context.Background(), fingerprint, u)
}

// GetUsers returns all users on the system that can have SSH keys.
func (s *sshKernelService) GetUsers() ([]user.User, error) {
	return s.sshAdapter.GetUsers(context.Background())
}

func (s *sshKernelService) GetUser(username string) (user.User, error) {
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

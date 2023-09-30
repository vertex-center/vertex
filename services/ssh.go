package services

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

type SSHService struct{}

func NewSSHService() SSHService {
	s := SSHService{}
	return s
}

func (s *SSHService) GetAll() ([]types.PublicKey, error) {
	var keys []types.PublicKey
	err := requests.URL(config.Current.HostKernel).
		Path("/api/security/ssh").
		ToJSON(&keys).
		Fetch(context.Background())
	return keys, err
}

func (s *SSHService) Add(key string) error {
	return requests.URL(config.Current.HostKernel).
		Path("/api/security/ssh").
		Post().
		BodyBytes([]byte(key)).
		Fetch(context.Background())
}

func (s *SSHService) Delete(fingerprint string) error {
	return requests.URL(config.Current.HostKernel).
		Pathf("/api/security/ssh/%s", fingerprint).
		Delete().
		Fetch(context.Background())
}

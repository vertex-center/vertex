package adapter

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

type SshKernelApiAdapter struct{}

func NewSshKernelApiAdapter() types.SshAdapterPort {
	return &SshKernelApiAdapter{}
}

func (a *SshKernelApiAdapter) GetAll() ([]types.PublicKey, error) {
	var keys []types.PublicKey
	err := requests.URL(config.Current.HostKernel).
		Path("/api/security/ssh").
		ToJSON(&keys).
		Fetch(context.Background())
	return keys, err
}

func (a *SshKernelApiAdapter) Add(key string) error {
	return requests.URL(config.Current.HostKernel).
		Path("/api/security/ssh").
		Post().
		BodyBytes([]byte(key)).
		Fetch(context.Background())
}

func (a *SshKernelApiAdapter) Remove(fingerprint string) error {
	return requests.URL(config.Current.HostKernel).
		Pathf("/api/security/ssh/%s", fingerprint).
		Delete().
		Fetch(context.Background())
}

package adapter

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

type SshKernelApiAdapter struct {
	config requests.Config
}

func NewSshKernelApiAdapter() types.SshAdapterPort {
	return &SshKernelApiAdapter{
		config: func(rb *requests.Builder) {
			rb.BaseURL(config.Current.HostKernel)
		},
	}
}

func (a *SshKernelApiAdapter) GetAll() ([]types.PublicKey, error) {
	var keys []types.PublicKey
	err := requests.New(a.config).
		Path("/api/security/ssh").
		ToJSON(&keys).
		Fetch(context.Background())
	return keys, err
}

func (a *SshKernelApiAdapter) Add(key string) error {
	return requests.New(a.config).
		Path("/api/security/ssh").
		Post().
		BodyBytes([]byte(key)).
		Fetch(context.Background())
}

func (a *SshKernelApiAdapter) Remove(fingerprint string) error {
	return requests.New(a.config).
		Pathf("/api/security/ssh/%s", fingerprint).
		Delete().
		Fetch(context.Background())
}

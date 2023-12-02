package port

import (
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	HardwareService interface {
		GetHost() (types.Host, error)
		GetCPUs() ([]types.CPU, error)
	}

	SshService interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string, username string) error
		Delete(fingerprint string, username string) error
		GetUsers() ([]string, error)
	}

	SshKernelService interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string, username string) error
		Delete(fingerprint string, username string) error
		GetUsers() ([]user.User, error)
	}
)

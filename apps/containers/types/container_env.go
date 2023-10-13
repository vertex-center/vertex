package types

import "github.com/google/uuid"

type ContainerEnvVariables map[string]string

type ContainerEnvAdapterPort interface {
	Save(uuid uuid.UUID, env ContainerEnvVariables) error
	Load(uuid uuid.UUID) (ContainerEnvVariables, error)
}

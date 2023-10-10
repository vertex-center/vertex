package types

import "github.com/google/uuid"

type InstanceEnvVariables map[string]string

type InstanceEnvAdapterPort interface {
	Save(uuid uuid.UUID, env InstanceEnvVariables) error
	Load(uuid uuid.UUID) (InstanceEnvVariables, error)
}

package types

import "github.com/google/uuid"

type InstanceServiceAdapterPort interface {
	Save(uuid uuid.UUID, service Service) error
	Load(uuid uuid.UUID) (Service, error)
	LoadRaw(uuid uuid.UUID) (interface{}, error)
}

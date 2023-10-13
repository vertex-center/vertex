package types

import "github.com/google/uuid"

type ContainerServiceAdapterPort interface {
	Save(uuid uuid.UUID, service Service) error
	Load(uuid uuid.UUID) (Service, error)
	LoadRaw(uuid uuid.UUID) (interface{}, error)
}

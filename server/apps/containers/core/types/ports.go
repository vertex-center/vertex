package types

import (
	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
)

var ErrInvalidPort = errors.NotValidf("port")

type (
	Ports []Port
	Port  struct {
		ID          uuid.UUID `json:"id"           db:"id"            example:"7e63ced7-4f4e-4b79-95ca-62930866f7bc"`
		ContainerID uuid.UUID `json:"container_id" db:"container_id"  example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		In          string    `json:"in"           db:"internal_port" example:"5432"` // Port in the container
		Out         string    `json:"out"          db:"external_port" example:"5432"` // Port exposed
	}

	PortFilters struct {
		ContainerID *uuid.UUID `json:"container_id,omitempty" db:"container_id"  example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
	}
)

func (p *Port) Validate() error {
	if len(p.In) == 0 || len(p.Out) == 0 {
		return ErrInvalidPort
	}
	return nil
}

package types

import "github.com/vertex-center/vertex/common/uuid"

type (
	Ports []Port
	Port  struct {
		ContainerID uuid.UUID `json:"container_id" db:"container_id"  example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		In          string    `json:"in"           db:"internal_port" example:"5432"`    // Port in the container
		Out         string    `json:"out"          db:"external_port" example:"DB_PORT"` // Port exposed
	}
)

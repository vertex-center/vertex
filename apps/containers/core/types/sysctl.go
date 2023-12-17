package types

import "github.com/vertex-center/uuid"

type (
	Sysctls []Sysctl
	Sysctl  struct {
		ContainerID uuid.UUID `json:"container_id" db:"container_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Name        string    `json:"name"         db:"name"`
		Value       string    `json:"value"        db:"value"`
	}
)

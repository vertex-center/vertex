package types

import "github.com/vertex-center/uuid"

type (
	Volumes []Volume
	Volume  struct {
		ContainerID uuid.UUID `json:"container_id" db:"container_id"   example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		In          string    `json:"in"           db:"internal_path"` // Path in the container
		Out         string    `json:"out"          db:"external_path"` // Path on the host
	}
)

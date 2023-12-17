package types

type (
	Capabilities []Capability
	Capability   struct {
		ContainerID ContainerID `json:"container_id" db:"container_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Name        string      `json:"name"         db:"name"`
	}
)

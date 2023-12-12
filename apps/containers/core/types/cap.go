package types

type (
	Capabilities []Capability
	Capability   struct {
		ContainerID ContainerID `json:"container_id" db:"container_id"`
		Name        string      `json:"name"         db:"name"`
	}
)

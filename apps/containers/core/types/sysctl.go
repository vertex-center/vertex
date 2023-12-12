package types

type (
	Sysctls []Sysctl
	Sysctl  struct {
		ContainerID ContainerID `json:"container_id" db:"container_id"`
		Name        string      `json:"name"         db:"name"`
		Value       string      `json:"value"        db:"value"`
	}
)

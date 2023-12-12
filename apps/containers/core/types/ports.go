package types

type (
	Ports []Port
	Port  struct {
		In  string `json:"in"  db:"internal_port" example:"5432"`    // Port in the container
		Out string `json:"out" db:"external_port" example:"DB_PORT"` // Port exposed
	}
)

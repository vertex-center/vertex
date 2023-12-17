package uuid

import "github.com/google/uuid"

var Nil UUID

// UUID is a wrapper around the google/uuid package
// to allow for custom marshalling and unmarshalling
type (
	UUID     struct{ uuid.UUID }
	NullUUID struct {
		uuid.NullUUID
		UUID UUID
	}
)

func New() UUID { return UUID{uuid.New()} }

func Parse(v string) (UUID, error) {
	u, err := uuid.Parse(v)
	return UUID{u}, err
}

func (UUID) Type() string   { return "string" }
func (UUID) Format() string { return "uuid" }

package uuid

import "github.com/google/uuid"

// UUID is a wrapper around the google/uuid package
// to allow for custom marshalling and unmarshalling
type UUID struct{ uuid.UUID }

func New() UUID { return UUID{uuid.New()} }

func (UUID) Type() string   { return "string" }
func (UUID) Format() string { return "uuid" }

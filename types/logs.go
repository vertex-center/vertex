package types

import "github.com/google/uuid"

const (
	LogKindOut       = "out"
	LogKindErr       = "err"
	LogKindVertexOut = "vertex_out"
	LogKindVertexErr = "vertex_err"
)

type LogLine struct {
	Id      int    `json:"id"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type InstanceLogsRepository interface {
	Open(uuid uuid.UUID) error
	Close(uuid uuid.UUID) error
	Push(uuid uuid.UUID, line LogLine)

	CloseAll() error
}

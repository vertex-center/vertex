package types

import "github.com/google/uuid"

const (
	LogKindOut       = "out"
	LogKindErr       = "err"
	LogKindDownload  = "download"
	LogKindVertexOut = "vertex_out"
	LogKindVertexErr = "vertex_err"
)

type LogLine struct {
	Id      int    `json:"id"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type InstanceLogsAdapterPort interface {
	Open(uuid uuid.UUID) error
	Close(uuid uuid.UUID) error
	Push(uuid uuid.UUID, line LogLine)

	// LoadBuffer will load the latest logs kept in memory.
	LoadBuffer(uuid uuid.UUID) ([]LogLine, error)

	CloseAll() error
}

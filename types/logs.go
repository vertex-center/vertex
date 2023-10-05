package types

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
)

const (
	LogKindOut       = "out"
	LogKindErr       = "err"
	LogKindDownload  = "download"
	LogKindDownloads = "downloads"
	LogKindVertexOut = "vertex_out"
	LogKindVertexErr = "vertex_err"
)

type InstanceLogsAdapterPort interface {
	Open(uuid uuid.UUID) error
	Close(uuid uuid.UUID) error

	Push(uuid uuid.UUID, line LogLine)
	Pop(uuid uuid.UUID) (LogLine, error)

	// LoadBuffer will load the latest logs kept in memory.
	LoadBuffer(uuid uuid.UUID) ([]LogLine, error)

	CloseAll() error
}

var ErrBufferEmpty = errors.New("the buffer is empty")

type LogLine struct {
	Id      int            `json:"id"`
	Kind    string         `json:"kind"`
	Message LogLineMessage `json:"message"`
}

type LogLineMessage interface{}

type LogLineMessageString struct {
	Value string `json:"value"`
}

func NewLogLineMessageString(s string) *LogLineMessageString {
	return &LogLineMessageString{
		Value: s,
	}
}

type LogLineMessageDownload struct {
	*DownloadProgress
}

func NewLogLineMessageDownload(p *DownloadProgress) *LogLineMessageDownload {
	return &LogLineMessageDownload{
		DownloadProgress: p,
	}
}

func (m *LogLineMessageDownload) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.DownloadProgress)
}

type DownloadProgressGroup []*DownloadProgress

type LogLineMessageDownloads struct {
	*DownloadProgressGroup
}

func NewLogLineMessageDownloads(p *DownloadProgress) *LogLineMessageDownloads {
	return &LogLineMessageDownloads{
		DownloadProgressGroup: &DownloadProgressGroup{p},
	}
}

func (m *LogLineMessageDownloads) Merge(progress *DownloadProgress) {
	if progress == nil {
		log.Error(errors.New("cannot merge nil progress group"))
		return
	}

	for i, p := range *m.DownloadProgressGroup {
		if p.ID == progress.ID {
			// This download ID already exists, update it.
			(*m.DownloadProgressGroup)[i] = progress
			return
		}
	}

	// This download ID does not exist, append it.
	*m.DownloadProgressGroup = append(*m.DownloadProgressGroup, progress)
}

func (m *LogLineMessageDownloads) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.DownloadProgressGroup)
}

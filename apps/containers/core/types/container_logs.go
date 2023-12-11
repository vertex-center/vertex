package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/common/log"
)

const (
	LogKindOut       = "out"
	LogKindErr       = "err"
	LogKindDownload  = "download"
	LogKindDownloads = "downloads"
	LogKindVertexOut = "vertex_out"
	LogKindVertexErr = "vertex_err"
)

var ErrBufferEmpty = errors.New("the buffer is empty")

type LogLine struct {
	Id      int            `json:"id"`
	Kind    string         `json:"kind"`
	Message LogLineMessage `json:"message"`
}

type LogLineMessage interface {
	String() string
}

type LogLineMessageString struct {
	Value string `json:"value"`
}

func NewLogLineMessageString(s string) *LogLineMessageString {
	return &LogLineMessageString{
		Value: s,
	}
}

func (m *LogLineMessageString) String() string {
	return m.Value
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

func (m *LogLineMessageDownload) String() string {
	if m.DownloadProgress == nil {
		return ""
	}
	return fmt.Sprintf("%v", *m.DownloadProgress)
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

func (m *LogLineMessageDownloads) String() string {
	if m.DownloadProgressGroup == nil {
		return ""
	}
	s := ""
	for _, p := range *m.DownloadProgressGroup {
		if p == nil {
			continue
		}
		s += fmt.Sprintf("%+v\n", *p)
	}
	s = s[:len(s)-1]
	return s
}

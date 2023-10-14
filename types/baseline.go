package types

import (
	"context"
	"errors"
	"reflect"
)

var ErrFailedToFetchBaseline = errors.New("failed to fetch baseline")

type Baseline struct {
	Date           string `json:"date"`            // Date of this release.
	Version        string `json:"version"`         // Public Version of the release.
	Vertex         string `json:"vertex"`          // Vertex version for this baseline Version.
	VertexClient   string `json:"vertex_client"`   // VertexClient version for this baseline Version.
	VertexServices string `json:"vertex_services"` // VertexServices version for this baseline Version.
}

func (b Baseline) GetVersionByID(id string) (string, error) {
	tp := reflect.TypeOf(b)
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		if field.Tag.Get("json") == id {
			value := reflect.ValueOf(b)
			return value.Field(i).String(), nil
		}
	}
	return "", errors.New("field not found")
}

type BaselinesAdapterPort interface {
	// GetLatest returns the latest available Baseline. This
	// will typically fetch the latest Baseline from a remote source.
	GetLatest(ctx context.Context, channel SettingsUpdatesChannel) (Baseline, error)
}

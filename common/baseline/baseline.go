package baseline

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

var ErrFailedToFetchBaseline = errors.New("failed to fetch baseline")

const (
	ChannelStable Channel = "stable"
	ChannelBeta   Channel = "beta"
)

type Channel string

type Baseline struct {
	Date           string `json:"date"`            // Date of this release.
	Version        string `json:"version"`         // Public Version of the release.
	Description    string `json:"description"`     // Condensed Description of the release.
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

func Fetch(ctx context.Context, version string) (Baseline, error) {
	var history []Baseline
	builder := requests.New().
		BaseURL("https://bl.vx.arra.red/").
		Pathf("versions.json").
		ToJSON(&history)

	url, err := builder.URL()
	if err != nil {
		return Baseline{}, err
	}

	log.Info("fetching baseline history", vlog.String("url", url.String()))

	err = builder.Fetch(ctx)
	if err != nil {
		return Baseline{}, fmt.Errorf("%w: %w", ErrFailedToFetchBaseline, err)
	}

	for _, baseline := range history {
		if baseline.Version == version {
			return baseline, nil
		}
	}

	return Baseline{}, fmt.Errorf("%w (%s)", ErrFailedToFetchBaseline, version)
}

func FetchLatest(ctx context.Context, channel Channel) (Baseline, error) {
	var baseline Baseline
	builder := requests.New().
		BaseURL("https://bl.vx.arra.red/").
		Pathf("%s.json", channel).
		ToJSON(&baseline)

	url, err := builder.URL()
	if err != nil {
		return baseline, err
	}

	log.Info("fetching latest baseline", vlog.String("url", url.String()))

	err = builder.Fetch(ctx)
	if err == nil {
		return baseline, nil
	}

	return baseline, fmt.Errorf("%w: %w", ErrFailedToFetchBaseline, err)
}

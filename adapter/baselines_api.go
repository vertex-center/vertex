package adapter

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type BaselinesApiAdapter struct {
	config requests.Config
}

func NewBaselinesApiAdapter() port.BaselinesAdapter {
	return &BaselinesApiAdapter{
		config: func(rb *requests.Builder) {
			rb.BaseURL("https://bl.vx.quentinguidee.dev/")
		},
	}
}

func (a *BaselinesApiAdapter) GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error) {
	var baseline types.Baseline
	builder := requests.New(a.config).
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

	return baseline, fmt.Errorf("%w: %w", types.ErrFailedToFetchBaseline, err)
}

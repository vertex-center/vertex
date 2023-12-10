package adapter

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

type baselinesApiAdapter struct {
	config requests.Config
}

func NewBaselinesApiAdapter() port.BaselinesAdapter {
	return &baselinesApiAdapter{
		config: func(rb *requests.Builder) {
			rb.BaseURL("https://bl.vx.arra.red/")
		},
	}
}

func (a *baselinesApiAdapter) GetLatest(ctx context.Context, channel types.UpdatesChannel) (types.Baseline, error) {
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

package adapter

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type BaselinesApiAdapter struct{}

func NewBaselinesApiAdapter() *BaselinesApiAdapter {
	return &BaselinesApiAdapter{}
}

func (a *BaselinesApiAdapter) GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error) {
	url := fmt.Sprintf("https://bl.vx.quentinguidee.dev/%s.json", channel)

	log.Info("fetching latest baseline", vlog.String("url", url))

	var baseline types.Baseline
	err := requests.URL(url).
		ToJSON(&baseline).
		Fetch(ctx)

	return baseline, fmt.Errorf("%w: %w", types.ErrFailedToFetchBaseline, err)
}

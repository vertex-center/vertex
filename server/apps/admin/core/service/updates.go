package service

import (
	"context"
	"fmt"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/baseline"
)

type updateService struct {
	ctx *apptypes.Context
}

func NewUpdateService(ctx *apptypes.Context) port.UpdateService {
	return &updateService{
		ctx: ctx,
	}
}

func (s *updateService) GetUpdate(channel baseline.Channel) (*types.Update, error) {
	bl, err := baseline.FetchLatest(context.Background(), channel)
	if err != nil {
		return nil, fmt.Errorf("fetch baseline: %w", err)
	}

	version := s.ctx.About().Version
	if bl.Version == version {
		return nil, nil
	}

	return &types.Update{
		Baseline: bl,
	}, nil
}

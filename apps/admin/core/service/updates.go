package service

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/baseline"
	"github.com/vertex-center/vertex/common/updater"
)

type updateService struct {
	uuid     uuid.UUID
	ctx      *apptypes.Context
	updaters []updater.Updater // updaters containers update logic for each dependency.
	updating atomic.Bool       // updating is true if an update is currently in progress.
}

func NewUpdateService(ctx *apptypes.Context, updaters []updater.Updater) port.UpdateService {
	return &updateService{
		uuid:     uuid.New(),
		ctx:      ctx,
		updaters: updaters,
	}
}

func (s *updateService) GetUpdate(channel baseline.Channel) (*types.Update, error) {
	bl, err := baseline.FetchLatest(context.Background(), channel)
	if err != nil {
		return nil, fmt.Errorf("fetch baseline: %w", err)
	}

	available, err := updater.CheckUpdates(bl, s.updaters...)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, nil
	}

	return &types.Update{
		Baseline: bl,
		Updating: s.updating.Load(),
	}, nil
}

func (s *updateService) InstallLatest(channel baseline.Channel) error {
	if !s.updating.CompareAndSwap(false, true) {
		return types.ErrAlreadyUpdating
	}
	defer s.updating.Store(false)

	bl, err := baseline.FetchLatest(context.Background(), channel)
	if err != nil {
		return fmt.Errorf("fetch baseline: %w", err)
	}

	return updater.Install(bl, s.updaters...)
}

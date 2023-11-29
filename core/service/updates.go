package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type UpdateService struct {
	uuid     uuid.UUID
	ctx      *types.VertexContext
	adapter  port.BaselinesAdapter
	updaters []types.Updater // updaters containers update logic for each dependency.
	updating atomic.Bool     // updating is true if an update is currently in progress.
}

func NewUpdateService(ctx *types.VertexContext, adapter port.BaselinesAdapter, updaters []types.Updater) port.UpdateService {
	s := &UpdateService{
		uuid:     uuid.New(),
		ctx:      ctx,
		adapter:  adapter,
		updaters: updaters,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *UpdateService) GetUpdate(channel types.UpdatesChannel) (*types.Update, error) {
	available := false
	update := types.Update{}

	latest, err := s.adapter.GetLatest(context.Background(), channel)
	if err != nil {
		return nil, err
	}

	log.Info("latest baseline fetched", vlog.Any("baseline", latest))

	for _, updater := range s.updaters {
		currentVersion, err := updater.CurrentVersion()
		if err != nil {
			return nil, err
		}

		latestVersion, err := latest.GetVersionByID(updater.ID())
		if err != nil {
			return nil, fmt.Errorf("'%w' when accessing '%s'", err, updater.ID())
		}

		if currentVersion != latestVersion {
			log.Info("update available",
				vlog.String("id", updater.ID()),
				vlog.String("current", currentVersion),
				vlog.String("latest", latestVersion))
			available = true
			update.Baseline = latest
		}
	}

	if !available {
		return nil, nil
	}

	update.Updating = s.updating.Load()

	return &update, nil
}

func (s *UpdateService) InstallLatest(channel types.UpdatesChannel) error {
	if !s.updating.CompareAndSwap(false, true) {
		return types.ErrAlreadyUpdating
	}
	defer s.updating.Store(false)

	latest, err := s.adapter.GetLatest(context.Background(), channel)
	if err != nil {
		return err
	}

	for _, updater := range s.updaters {
		v, err := latest.GetVersionByID(updater.ID())
		if err != nil {
			return err
		}

		err = updater.Install(v)
		if err != nil {
			return err
		}
	}

	s.ctx.DispatchEvent(types.EventVertexUpdated{})
	return nil
}

func (s *UpdateService) firstSetup() error {
	var missingDeps []types.Updater
	for _, updater := range s.updaters {
		if !updater.IsInstalled() {
			missingDeps = append(missingDeps, updater)
		}
	}

	if len(missingDeps) == 0 {
		log.Info("all dependencies are already installed")
		return nil
	}

	log.Info("installing missing dependencies", vlog.Any("count", len(missingDeps)))

	latest, err := s.adapter.GetLatest(context.Background(), types.UpdatesChannelStable)
	if err != nil {
		return err
	}

	for _, updater := range missingDeps {
		version, err := latest.GetVersionByID(updater.ID())
		if err != nil {
			return err
		}

		err = updater.Install(version)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UpdateService) OnEvent(e event.Event) error {
	switch e.(type) {
	case types.EventServerStart:
		err := s.firstSetup()
		if err != nil {
			log.Error(err)
			log.Error(errors.New("failed to fetch latest baseline. panic because vertex cannot run without missing dependencies"))
			os.Exit(1)
		}

		err = config.Current.Apply()
		if err != nil {
			log.Error(fmt.Errorf("failed to apply the current configuration: %w", err))
			os.Exit(1)
		}
	}
	return nil
}

func (s *UpdateService) GetUUID() uuid.UUID {
	return s.uuid
}

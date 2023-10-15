package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
	"os"
	"sync"
)

type UpdateService struct {
	uuid     uuid.UUID
	ctx      *types.VertexContext
	adapter  types.BaselinesAdapterPort
	updaters []types.Updater // updaters containers update logic for each dependency.
	updateMu sync.Mutex      // updateMu is used to prevent multiple updates at the same time.
}

func NewUpdateService(ctx *types.VertexContext, adapter types.BaselinesAdapterPort, updaters []types.Updater) *UpdateService {
	s := &UpdateService{
		uuid:     uuid.New(),
		ctx:      ctx,
		adapter:  adapter,
		updaters: updaters,
		updateMu: sync.Mutex{},
	}
	s.ctx.AddListener(s)
	return s
}

func (s *UpdateService) GetUpdate(channel types.SettingsUpdatesChannel) (*types.Update, error) {
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
			available = true
			update.Baseline = latest
		}
	}

	if !available {
		return nil, nil
	}

	return &update, nil
}

func (s *UpdateService) InstallLatest(channel types.SettingsUpdatesChannel) error {
	if !s.updateMu.TryLock() {
		return types.ErrAlreadyUpdating
	}
	defer s.updateMu.Unlock()

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

	latest, err := s.adapter.GetLatest(context.Background(), types.SettingsUpdatesChannelStable)
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

func (s *UpdateService) OnEvent(e interface{}) {
	switch e.(type) {
	case types.EventServerStart:
		err := s.firstSetup()
		if err != nil {
			log.Error(err)
			err = errors.New("failed to fetch latest baseline. panic because vertex cannot run without missing dependencies")
			log.Error(err)
			os.Exit(1)
		}

		err = config.Current.Apply()
		if err != nil {
			log.Error(fmt.Errorf("failed to apply the current configuration: %v", err))
			os.Exit(1)
		}
	}
}

func (s *UpdateService) GetUUID() uuid.UUID {
	return s.uuid
}
package services

import (
	"bufio"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type InstanceRunnerService struct {
	eventsAdapter types.EventAdapterPort
	adapter       types.RunnerAdapterPort
}

func NewInstanceRunnerService(eventsAdapter types.EventAdapterPort, dockerRunnerAdapter types.RunnerAdapterPort) InstanceRunnerService {
	return InstanceRunnerService{
		eventsAdapter: eventsAdapter,
		adapter:       dockerRunnerAdapter,
	}
}

func (s *InstanceRunnerService) Install(uuid uuid.UUID, service types.Service) error {
	if service.Methods.Docker == nil {
		return ErrInstallMethodDoesNotExists
	}

	dir := path.Join(storage.Path, uuid.String())
	if service.Methods.Docker.Clone != nil {
		err := storage.CloneRepository(dir, service.Methods.Docker.Clone.Repository)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *InstanceRunnerService) Delete(inst *types.Instance) error {
	return s.adapter.Delete(inst)
}

// Start starts an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is already running, it returns ErrInstanceAlreadyRunning.
func (s *InstanceRunnerService) Start(inst *types.Instance) error {
	if inst.IsBusy() {
		return nil
	}

	s.eventsAdapter.Send(types.EventInstanceLog{
		InstanceUUID: inst.UUID,
		Kind:         types.LogKindOut,
		Message:      "Starting instance...",
	})

	log.Info("starting instance",
		vlog.String("uuid", inst.UUID.String()),
	)

	if inst.IsRunning() {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: inst.UUID,
			Kind:         types.LogKindVertexErr,
			Message:      ErrInstanceAlreadyRunning.Error(),
		})
		return ErrInstanceAlreadyRunning
	}

	setStatus := func(status string) {
		s.setStatus(inst, status)
	}

	var runner types.RunnerAdapterPort
	if inst.IsDockerized() {
		runner = s.adapter
	} else {
		return fmt.Errorf("instance is not dockerized")
	}

	stdout, stderr, err := runner.Start(inst, setStatus)
	if err != nil {
		s.setStatus(inst, types.InstanceStatusError)
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}

			if strings.HasPrefix(scanner.Text(), "DOWNLOAD") {
				s.eventsAdapter.Send(types.EventInstanceLog{
					InstanceUUID: inst.UUID,
					Kind:         types.LogKindDownload,
					Message:      strings.TrimPrefix(scanner.Text(), "DOWNLOAD"),
				})
				continue
			}

			s.eventsAdapter.Send(types.EventInstanceLog{
				InstanceUUID: inst.UUID,
				Kind:         types.LogKindOut,
				Message:      scanner.Text(),
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			s.eventsAdapter.Send(types.EventInstanceLog{
				InstanceUUID: inst.UUID,
				Kind:         types.LogKindErr,
				Message:      scanner.Text(),
			})
		}
	}()

	// Wait for the instance until stopped
	wg.Wait()

	// Log stopped
	s.eventsAdapter.Send(types.EventInstanceLog{
		InstanceUUID: inst.UUID,
		Kind:         types.LogKindVertexOut,
		Message:      "Stopping instance...",
	})
	log.Info("stopping instance",
		vlog.String("uuid", inst.UUID.String()),
	)

	return nil
}

// Stop stops an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is not running, it returns ErrInstanceNotRunning.
func (s *InstanceRunnerService) Stop(inst *types.Instance) error {
	if inst.IsBusy() {
		return nil
	}

	if !inst.IsRunning() {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: inst.UUID,
			Kind:         types.LogKindVertexErr,
			Message:      ErrInstanceNotRunning.Error(),
		})
		return ErrInstanceNotRunning
	}

	s.setStatus(inst, types.InstanceStatusStopping)

	var err error
	if inst.IsDockerized() {
		err = s.adapter.Stop(inst)
	} else {
		return fmt.Errorf("inst is not dockerized")
	}

	if err == nil {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: inst.UUID,
			Kind:         types.LogKindVertexOut,
			Message:      "Instance stopped.",
		})

		log.Info("inst stopped",
			vlog.String("uuid", inst.UUID.String()),
		)

		s.setStatus(inst, types.InstanceStatusOff)
	} else {
		s.setStatus(inst, types.InstanceStatusRunning)
	}

	return err
}

func (s *InstanceRunnerService) GetDockerContainerInfo(inst types.Instance) (map[string]any, error) {
	if !inst.IsDockerized() {
		return nil, errors.New("instance is not using docker")
	}

	info, err := s.adapter.Info(inst)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *InstanceRunnerService) GetAllVersions(inst *types.Instance, useCache bool) ([]string, error) {
	if !useCache || len(inst.CacheVersions) == 0 {
		versions, err := s.adapter.GetAllVersions(*inst)
		if err != nil {
			return nil, err
		}
		inst.CacheVersions = versions
	}

	return inst.CacheVersions, nil
}

func (s *InstanceRunnerService) CheckForUpdates(inst *types.Instance) error {
	return s.adapter.CheckForUpdates(inst)
}

// RecreateContainer recreates a container by its UUID.
func (s *InstanceRunnerService) RecreateContainer(inst *types.Instance) error {
	if !inst.IsDockerized() {
		return nil
	}

	if inst.IsRunning() {
		err := s.adapter.Stop(inst)
		if err != nil {
			return err
		}
	}

	err := s.adapter.Delete(inst)
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	go func() {
		err := s.Start(inst)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	return nil
}

func (s *InstanceRunnerService) setStatus(inst *types.Instance, status string) {
	if inst.Status == status {
		return
	}

	var name string
	if inst.DisplayName == nil {
		name = inst.Service.Name
	} else {
		name = *inst.DisplayName
	}

	inst.Status = status
	s.eventsAdapter.Send(types.EventInstancesChange{})
	s.eventsAdapter.Send(types.EventInstanceStatusChange{
		InstanceUUID: inst.UUID,
		Name:         name,
		Status:       status,
	})
}

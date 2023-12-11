package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"path"
	"strings"
	"sync"

	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/pkg/event"
	vstorage "github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
)

type containerRunnerService struct {
	ctx     *app.Context
	adapter port.RunnerAdapter
}

func NewRunnerService(ctx *app.Context, adapter port.RunnerAdapter) port.RunnerService {
	return &containerRunnerService{
		ctx:     ctx,
		adapter: adapter,
	}
}

func (s *containerRunnerService) Install(ctx context.Context, uuid types.ContainerID, service types.Service) error {
	if service.Methods.Docker == nil {
		return ErrInstallMethodDoesNotExists
	}

	dir := path.Join(storage.FSPath, uuid.String())
	if service.Methods.Docker.Clone != nil {
		err := vstorage.CloneRepository(dir, service.Methods.Docker.Clone.Repository)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *containerRunnerService) Delete(ctx context.Context, inst *types.Container) error {
	err := s.adapter.DeleteMounts(ctx, inst)
	if err != nil {
		return err
	}
	return s.adapter.DeleteContainer(ctx, inst)
}

// Start starts a container by its UUID.
// If the container does not exist, it returns ErrContainerNotFound.
// If the container is already running, it returns ErrContainerAlreadyRunning.
func (s *containerRunnerService) Start(ctx context.Context, inst *types.Container) error {
	if inst.IsBusy() {
		return nil
	}

	s.ctx.DispatchEvent(types.EventContainerLog{
		ContainerUUID: inst.UUID,
		Kind:          types.LogKindOut,
		Message:       types.NewLogLineMessageString("Starting container..."),
	})

	log.Info("starting container",
		vlog.String("uuid", inst.UUID.String()),
	)

	if inst.IsRunning() {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerUUID: inst.UUID,
			Kind:          types.LogKindVertexErr,
			Message:       types.NewLogLineMessageString(ErrContainerAlreadyRunning.Error()),
		})
		return ErrContainerAlreadyRunning
	}

	setStatus := func(status string) {
		s.setStatus(inst, status)
	}

	stdout, stderr, err := s.adapter.Start(ctx, inst, setStatus)
	if err != nil {
		s.setStatus(inst, types.ContainerStatusError)
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}

			if strings.HasPrefix(scanner.Text(), "DOWNLOAD") {
				msg := strings.TrimPrefix(scanner.Text(), "DOWNLOAD")

				var downloadProgress types.DownloadProgress
				err := json.Unmarshal([]byte(msg), &downloadProgress)
				if err != nil {
					log.Error(err)
					continue
				}

				s.ctx.DispatchEvent(types.EventContainerLog{
					ContainerUUID: inst.UUID,
					Kind:          types.LogKindDownload,
					Message:       types.NewLogLineMessageDownload(&downloadProgress),
				})
				continue
			}

			s.ctx.DispatchEvent(types.EventContainerLog{
				ContainerUUID: inst.UUID,
				Kind:          types.LogKindOut,
				Message:       types.NewLogLineMessageString(scanner.Text()),
			})
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			s.ctx.DispatchEvent(types.EventContainerLog{
				ContainerUUID: inst.UUID,
				Kind:          types.LogKindErr,
				Message:       types.NewLogLineMessageString(scanner.Text()),
			})
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := s.WaitStatus(ctx, inst, types.ContainerStatusRunning)
		if err != nil {
			log.Error(err)
		}
	}()
	wg.Wait()

	return nil
}

// Stop stops a container by its UUID.
// If the container does not exist, it returns ErrContainerNotFound.
// If the container is not running, it returns ErrContainerNotRunning.
func (s *containerRunnerService) Stop(ctx context.Context, inst *types.Container) error {
	if inst.IsBusy() {
		return nil
	}

	if !inst.IsRunning() {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerUUID: inst.UUID,
			Kind:          types.LogKindVertexErr,
			Message:       types.NewLogLineMessageString(ErrContainerNotRunning.Error()),
		})
		return ErrContainerNotRunning
	}

	// Log stopped
	s.ctx.DispatchEvent(types.EventContainerLog{
		ContainerUUID: inst.UUID,
		Kind:          types.LogKindVertexOut,
		Message:       types.NewLogLineMessageString("Stopping container..."),
	})
	log.Info("stopping container",
		vlog.String("uuid", inst.UUID.String()),
	)

	s.setStatus(inst, types.ContainerStatusStopping)

	err := s.adapter.Stop(ctx, inst)
	if err == nil {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerUUID: inst.UUID,
			Kind:          types.LogKindVertexOut,
			Message:       types.NewLogLineMessageString("Container stopped."),
		})

		log.Info("container stopped",
			vlog.String("uuid", inst.UUID.String()),
		)

		s.setStatus(inst, types.ContainerStatusOff)
	} else {
		s.setStatus(inst, types.ContainerStatusRunning)
	}

	return err
}

func (s *containerRunnerService) GetDockerContainerInfo(ctx context.Context, inst types.Container) (map[string]any, error) {
	return s.adapter.Info(ctx, inst)
}

func (s *containerRunnerService) GetAllVersions(ctx context.Context, inst *types.Container, useCache bool) ([]string, error) {
	if !useCache || len(inst.CacheVersions) == 0 {
		versions, err := s.adapter.GetAllVersions(ctx, *inst)
		if err != nil {
			return nil, err
		}
		inst.CacheVersions = versions
	}

	return inst.CacheVersions, nil
}

func (s *containerRunnerService) CheckForUpdates(ctx context.Context, inst *types.Container) error {
	return s.adapter.CheckForUpdates(ctx, inst)
}

// RecreateContainer recreates a container by its UUID.
func (s *containerRunnerService) RecreateContainer(ctx context.Context, inst *types.Container) error {
	if inst.IsRunning() {
		err := s.adapter.Stop(ctx, inst)
		if err != nil {
			return err
		}
	}

	err := s.adapter.DeleteContainer(ctx, inst)
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	return s.Start(ctx, inst)
}

func (s *containerRunnerService) WaitStatus(ctx context.Context, inst *types.Container, status string) error {
	statusChan := make(chan string)
	defer close(statusChan)

	if inst.Status == status {
		return nil
	}

	l := event.NewTempListener(func(e event.Event) error {
		switch e := e.(type) {
		case types.EventContainerStatusChange:
			if e.ContainerUUID != inst.UUID {
				return nil
			}
			statusChan <- e.Status
		}
		return nil
	})

	s.ctx.AddListener(l)
	defer s.ctx.RemoveListener(l)

	for e := range statusChan {
		if e == status {
			return nil
		}
	}

	return errors.New("wait status timeout")
}

func (s *containerRunnerService) setStatus(inst *types.Container, status string) {
	if inst.Status == status {
		return
	}

	inst.Status = status
	s.ctx.DispatchEvent(types.EventContainersChange{})
	s.ctx.DispatchEvent(types.EventContainerStatusChange{
		ContainerUUID: inst.UUID,
		ServiceID:     inst.Service.ID,
		Container:     *inst,
		Name:          inst.DisplayName,
		Status:        status,
	})
}

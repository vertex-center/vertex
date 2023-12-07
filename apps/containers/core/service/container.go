package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vlog"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
)

var (
	ErrContainerAlreadyExists     = errors.New("container already exists")
	ErrContainerAlreadyRunning    = errors.New("the container is already running")
	ErrContainerNotRunning        = errors.New("the container is not running")
	ErrInstallMethodDoesNotExists = errors.New("this install method doesn't exist for this service")
)

type containerService struct {
	uuid uuid.UUID
	ctx  *app.Context

	containerAdapter port.ContainerAdapter

	containerRunnerService   port.ContainerRunnerService
	containerServiceService  port.ContainerServiceService
	containerEnvService      port.ContainerEnvService
	containerSettingsService port.ContainerSettingsService
	serviceService           port.ServiceService

	containers      map[uuid.UUID]*types.Container
	containersMutex *sync.RWMutex
}

type ContainerServiceParams struct {
	Ctx *app.Context

	ContainerAdapter port.ContainerAdapter

	ContainerRunnerService   port.ContainerRunnerService
	ContainerServiceService  port.ContainerServiceService
	ContainerEnvService      port.ContainerEnvService
	ContainerSettingsService port.ContainerSettingsService
	ServiceService           port.ServiceService
}

func NewContainerService(params ContainerServiceParams) port.ContainerService {
	s := &containerService{
		uuid: uuid.New(),
		ctx:  params.Ctx,

		containerAdapter: params.ContainerAdapter,

		containerRunnerService:   params.ContainerRunnerService,
		containerServiceService:  params.ContainerServiceService,
		containerEnvService:      params.ContainerEnvService,
		containerSettingsService: params.ContainerSettingsService,
		serviceService:           params.ServiceService,

		containers:      make(map[uuid.UUID]*types.Container),
		containersMutex: &sync.RWMutex{},
	}

	s.ctx.AddListener(s)

	return s
}

func (s *containerService) Get(uuid uuid.UUID) (*types.Container, error) {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	container, ok := s.containers[uuid]
	if !ok {
		return nil, types.ErrContainerNotFound
	}
	return container, nil
}

func (s *containerService) GetAll() map[uuid.UUID]*types.Container {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	return s.containers
}

func (s *containerService) GetTags() []string {
	var tags []string

	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	for _, inst := range s.containers {
		for _, tag := range inst.Tags {
			found := false
			for _, t := range tags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

// Search returns all containers that match the query.
func (s *containerService) Search(query types.ContainerSearchQuery) map[uuid.UUID]*types.Container {
	containers := map[uuid.UUID]*types.Container{}

	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	for _, inst := range s.containers {
		if query.Features != nil {
			if !inst.HasFeatureIn(*query.Features) {
				continue
			}
		}
		if query.Tags != nil {
			if !inst.HasTagIn(*query.Tags) {
				continue
			}
		}
		containers[inst.UUID] = inst
	}

	return containers
}

func (s *containerService) Exists(uuid uuid.UUID) bool {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	return s.containers[uuid] != nil
}

// Delete deletes a container by its UUID.
// If the container is still running, it returns ErrContainerStillRunning.
func (s *containerService) Delete(inst *types.Container) error {
	serviceID := inst.Service.ID

	if inst.IsRunning() {
		return types.ErrContainerStillRunning
	}

	err := s.containerRunnerService.Delete(inst)
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	err = s.containerAdapter.Delete(inst.UUID)
	if err != nil {
		return err
	}

	s.containersMutex.Lock()
	defer s.containersMutex.Unlock()
	delete(s.containers, inst.UUID)

	s.ctx.DispatchEvent(types.EventContainerDeleted{
		ContainerUUID: inst.UUID,
		ServiceID:     serviceID,
	})

	s.ctx.DispatchEvent(types.EventContainersChange{})

	return nil
}

func (s *containerService) StartAll() {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	var ids []uuid.UUID

	for _, inst := range s.containers {
		// vertex containers autostart are managed by the startup service.
		if inst.LaunchOnStartup() {
			ids = append(ids, inst.UUID)
		}
	}

	if len(ids) == 0 {
		return
	}

	log.Info("trying to ping Google...")

	// Wait for internet connection
	timeout, cancelTimeout := context.WithTimeout(context.Background(), 60*time.Second)
	err := net.WaitInternetConn(timeout)
	cancelTimeout()
	if err != nil {
		log.Error(err)
		return
	}

	// Start them
	for _, id := range ids {
		go func(id uuid.UUID) {
			inst, err := s.Get(id)
			if err != nil {
				log.Error(err)
				return
			}

			err = s.containerRunnerService.Start(inst)
			if err != nil {
				log.Warn("failed to auto-start the container",
					vlog.String("uuid", inst.UUID.String()),
					vlog.String("reason", err.Error()),
				)
			}
		}(id)
	}
}

func (s *containerService) StopAll() {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	for _, inst := range s.containers {
		err := s.containerRunnerService.Stop(inst)
		if err != nil {
			log.Error(err)
		}
	}

	s.ctx.DispatchEvent(types.EventContainersStopped{})
}

func (s *containerService) LoadAll() {
	uuids, err := s.containerAdapter.GetAll()
	if err != nil {
		return
	}

	loaded := 0
	for _, id := range uuids {
		err := s.load(id)
		if err != nil {
			log.Error(err)
			continue
		}
		loaded += 1
	}

	s.ctx.DispatchEvent(types.EventContainersLoaded{
		Count: loaded,
	})
}

func (s *containerService) DeleteAll() {
	all := s.GetAll()
	for _, inst := range all {
		err := s.Delete(inst)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *containerService) Install(service types.Service, method string) (*types.Container, error) {
	id := uuid.New()
	err := s.containerAdapter.Create(id)
	if err != nil {
		return nil, err
	}

	err = s.containerRunnerService.Install(id, service)
	if err != nil {
		return nil, err
	}

	tempContainer := &types.Container{
		UUID:    id,
		Service: service,
	}

	err = s.containerServiceService.Save(tempContainer, service)
	if err != nil {
		return nil, err
	}

	err = s.load(id)
	if err != nil {
		return nil, err
	}

	inst, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	inst.ContainerSettings.InstallMethod = &method
	err = s.containerSettingsService.Save(inst, inst.ContainerSettings)
	if err != nil {
		return nil, err
	}

	inst.ResetDefaultEnv()
	err = s.containerEnvService.Save(inst, inst.Env)
	if err != nil {
		return nil, err
	}

	s.ctx.DispatchEvent(types.EventContainerCreated{})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return inst, nil
}

func (s *containerService) CheckForUpdates() (map[uuid.UUID]*types.Container, error) {
	for _, inst := range s.GetAll() {
		err := s.containerRunnerService.CheckForUpdates(inst)
		if err != nil {
			return s.GetAll(), err
		}
	}

	return s.GetAll(), nil
}

func (s *containerService) load(uuid uuid.UUID) error {
	service, err := s.containerServiceService.Load(uuid)
	if err != nil {
		return err
	}

	inst := types.NewContainer(uuid, service)

	err = s.containerSettingsService.Load(&inst)
	if err != nil {
		return err
	}

	err = s.containerEnvService.Load(&inst)
	if err != nil {
		return err
	}

	latestService, err := s.serviceService.GetById(service.ID)
	if err != nil {
		log.Error(err)
	} else {
		err = s.containerServiceService.CheckForUpdate(&inst, latestService)
		if err != nil {
			log.Error(err)
		}
	}

	if !s.Exists(uuid) {
		s.containersMutex.Lock()
		defer s.containersMutex.Unlock()
		s.containers[uuid] = &inst
	} else {
		return ErrContainerAlreadyExists
	}

	s.ctx.DispatchEvent(types.EventContainerLoaded{
		Container: &inst,
	})

	return nil
}

func (s *containerService) SetDatabases(inst *types.Container, databases map[string]uuid.UUID, options map[string]*types.SetDatabasesOptions) error {
	for db := range databases {
		if _, ok := inst.Service.Databases[db]; !ok {
			return types.ErrDatabaseIDNotFound
		}
	}

	inst.Databases = databases
	err := s.containerSettingsService.Save(inst, inst.ContainerSettings)
	if err != nil {
		return err
	}
	return s.remapDatabaseEnv(inst, options)
}

// remapDatabaseEnv remaps the environment variables of a container.
func (s *containerService) remapDatabaseEnv(inst *types.Container, options map[string]*types.SetDatabasesOptions) error {
	for databaseID, databaseContainerUUID := range inst.Databases {
		db, err := s.Get(databaseContainerUUID)
		if err != nil {
			return err
		}

		host := config.Current.URL("vertex").String()

		dbEnvNames := (*db.Service.Features.Databases)[0]
		iEnvNames := inst.Service.Databases[databaseID].Names

		inst.Env[iEnvNames.Host] = host
		inst.Env[iEnvNames.Port] = db.Env[dbEnvNames.Port]
		if dbEnvNames.Username != nil {
			inst.Env[iEnvNames.Username] = db.Env[*dbEnvNames.Username]
		}
		if dbEnvNames.Password != nil {
			inst.Env[iEnvNames.Password] = db.Env[*dbEnvNames.Password]
		}

		if options != nil {
			if modifiedFeature, ok := options[databaseID]; ok {
				if modifiedFeature != nil && modifiedFeature.DatabaseName != nil {
					inst.Env[iEnvNames.Database] = *modifiedFeature.DatabaseName
					continue
				}
			}
		}

		if dbEnvNames.DefaultDatabase != nil {
			inst.Env[iEnvNames.Database] = db.Env[*dbEnvNames.DefaultDatabase]
			continue
		}

		delete(inst.Env, iEnvNames.Database)

	}

	return s.containerEnvService.Save(inst, inst.Env)
}

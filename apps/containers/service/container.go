package service

import (
	"errors"
	"github.com/vertex-center/vertex/core/types/app"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/types"
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

type ContainerService struct {
	uuid uuid.UUID
	ctx  *app.Context

	containerAdapter types.ContainerAdapterPort

	containerRunnerService   *ContainerRunnerService
	containerServiceService  *ContainerServiceService
	containerEnvService      *ContainerEnvService
	containerSettingsService *ContainerSettingsService

	containers      map[uuid.UUID]*types.Container
	containersMutex *sync.RWMutex
}

type ContainerServiceParams struct {
	Ctx *app.Context

	ContainerAdapter types.ContainerAdapterPort

	ContainerRunnerService   *ContainerRunnerService
	ContainerServiceService  *ContainerServiceService
	ContainerEnvService      *ContainerEnvService
	ContainerSettingsService *ContainerSettingsService
}

func NewContainerService(params ContainerServiceParams) *ContainerService {
	s := &ContainerService{
		uuid: uuid.New(),
		ctx:  params.Ctx,

		containerAdapter: params.ContainerAdapter,

		containerRunnerService:   params.ContainerRunnerService,
		containerServiceService:  params.ContainerServiceService,
		containerEnvService:      params.ContainerEnvService,
		containerSettingsService: params.ContainerSettingsService,

		containers:      make(map[uuid.UUID]*types.Container),
		containersMutex: &sync.RWMutex{},
	}

	s.ctx.AddListener(s)

	return s
}

// Get returns an container by its UUID. If the container does not exist,
// it returns ErrContainerNotFound.
func (s *ContainerService) Get(uuid uuid.UUID) (*types.Container, error) {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	container, ok := s.containers[uuid]
	if !ok {
		return nil, types.ErrContainerNotFound
	}
	return container, nil
}

func (s *ContainerService) GetAll() map[uuid.UUID]*types.Container {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	return s.containers
}

func (s *ContainerService) GetTags() []string {
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
func (s *ContainerService) Search(query types.ContainerSearchQuery) map[uuid.UUID]*types.Container {
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

func (s *ContainerService) Exists(uuid uuid.UUID) bool {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	return s.containers[uuid] != nil
}

// Delete deletes an container by its UUID.
// If the container is still running, it returns ErrContainerStillRunning.
func (s *ContainerService) Delete(inst *types.Container) error {
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

func (s *ContainerService) StartAll() {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	var ids []uuid.UUID

	for _, inst := range s.containers {
		// vertex containers autostart are managed by the startup service.
		if inst.LaunchOnStartup() && !inst.HasTag("vertex") {
			ids = append(ids, inst.UUID)
		}
	}

	if len(ids) == 0 {
		return
	}

	log.Info("trying to ping Google...")

	// Wait for internet connection
	err := net.Wait("google.com:80")
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
				log.Error(err)
			}
		}(id)
	}
}

func (s *ContainerService) StopAll() {
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

func (s *ContainerService) Install(service types.Service, method string) (*types.Container, error) {
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

func (s *ContainerService) CheckForUpdates() (map[uuid.UUID]*types.Container, error) {
	for _, inst := range s.GetAll() {
		err := s.containerRunnerService.CheckForUpdates(inst)
		if err != nil {
			return s.GetAll(), err
		}
	}

	return s.GetAll(), nil
}

func (s *ContainerService) LoadAll() {
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

func (s *ContainerService) DeleteAll() {
	all := s.GetAll()
	for _, inst := range all {
		err := s.Delete(inst)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *ContainerService) load(uuid uuid.UUID) error {
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

	err = s.containerServiceService.CheckForUpdate(&inst, service)
	if err != nil {
		return err
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

func (s *ContainerService) SetDatabases(inst *types.Container, databases map[string]uuid.UUID) error {
	inst.Databases = databases
	err := s.containerSettingsService.Save(inst, inst.ContainerSettings)
	if err != nil {
		return err
	}
	return s.remapDatabaseEnv(inst)
}

// remapDatabaseEnv remaps the environment variables of an container.
func (s *ContainerService) remapDatabaseEnv(inst *types.Container) error {
	for databaseID, databaseContainerUUID := range inst.Databases {
		db, err := s.Get(databaseContainerUUID)
		if err != nil {
			return err
		}

		host := config.Current.Host

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
	}

	return s.containerEnvService.Save(inst, inst.Env)
}

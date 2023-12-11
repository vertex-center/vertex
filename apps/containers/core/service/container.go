package service

import (
	"context"
	"sync"
	"time"

	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
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

	runnerService           port.RunnerService
	containerServiceService port.ContainerServiceService
	envService              port.EnvService
	settingsService         port.SettingsService
	serviceService          port.ServiceService

	containers      map[types.ContainerID]*types.Container
	containersMutex *sync.RWMutex
}

type ContainerServiceParams struct {
	Ctx *app.Context

	ContainerAdapter port.ContainerAdapter

	RunnerService           port.RunnerService
	ContainerServiceService port.ContainerServiceService
	EnvService              port.EnvService
	SettingsService         port.SettingsService
	ServiceService          port.ServiceService
}

func NewContainerService(params ContainerServiceParams) port.ContainerService {
	s := &containerService{
		uuid: uuid.New(),
		ctx:  params.Ctx,

		containerAdapter: params.ContainerAdapter,

		runnerService:           params.RunnerService,
		containerServiceService: params.ContainerServiceService,
		envService:              params.EnvService,
		settingsService:         params.SettingsService,
		serviceService:          params.ServiceService,

		containers:      make(map[types.ContainerID]*types.Container),
		containersMutex: &sync.RWMutex{},
	}

	s.ctx.AddListener(s)

	return s
}

// Get returns a container by its UUID.
// If the container doesn't exist, it returns ErrContainerNotFound.
func (s *containerService) Get(ctx context.Context, uuid types.ContainerID) (*types.Container, error) {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	container, ok := s.containers[uuid]
	if !ok {
		return nil, types.ErrContainerNotFound
	}
	return container, nil
}

func (s *containerService) GetAll(ctx context.Context) map[types.ContainerID]*types.Container {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	return s.containers
}

func (s *containerService) GetTags(ctx context.Context) []string {
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
func (s *containerService) Search(ctx context.Context, query types.ContainerSearchQuery) map[types.ContainerID]*types.Container {
	containers := map[types.ContainerID]*types.Container{}

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

func (s *containerService) Exists(ctx context.Context, uuid types.ContainerID) bool {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	return s.containers[uuid] != nil
}

func (s *containerService) Delete(ctx context.Context, uuid types.ContainerID) error {
	inst, err := s.Get(ctx, uuid)
	if err != nil {
		return err
	}

	if inst.IsRunning() {
		return types.ErrContainerStillRunning
	}

	err = s.runnerService.Delete(ctx, inst)
	if err != nil && !errors.Is(err, errors.NotFound) {
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
		ServiceID:     inst.Service.ID,
	})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return nil
}

func (s *containerService) StartAll(ctx context.Context) {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	var ids []types.ContainerID

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
		go func(id types.ContainerID) {
			inst, err := s.Get(ctx, id)
			if err != nil {
				log.Error(err)
				return
			}

			err = s.runnerService.Start(ctx, inst)
			if err != nil {
				log.Warn("failed to auto-start the container",
					vlog.String("uuid", inst.UUID.String()),
					vlog.String("reason", err.Error()),
				)
			}
		}(id)
	}
}

func (s *containerService) StopAll(ctx context.Context) {
	s.containersMutex.RLock()
	defer s.containersMutex.RUnlock()

	for _, inst := range s.containers {
		err := s.runnerService.Stop(ctx, inst)
		if err != nil {
			log.Error(err)
		}
	}

	s.ctx.DispatchEvent(types.EventContainersStopped{})
}

func (s *containerService) LoadAll(ctx context.Context) {
	uuids, err := s.containerAdapter.GetAll()
	if err != nil {
		return
	}

	loaded := 0
	for _, id := range uuids {
		err := s.load(ctx, id)
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

func (s *containerService) DeleteAll(ctx context.Context) {
	all := s.GetAll(ctx)
	for _, inst := range all {
		err := s.Delete(ctx, inst.UUID)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *containerService) Install(ctx context.Context, service types.Service, method string) (*types.Container, error) {
	id := types.NewContainerID()
	err := s.containerAdapter.Create(id)
	if err != nil {
		return nil, err
	}

	err = s.runnerService.Install(ctx, id, service)
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

	err = s.load(ctx, id)
	if err != nil {
		return nil, err
	}

	inst, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	inst.ContainerSettings.InstallMethod = &method
	err = s.settingsService.Save(inst, inst.ContainerSettings)
	if err != nil {
		return nil, err
	}

	inst.ResetDefaultEnv()
	err = s.envService.Save(inst, inst.Env)
	if err != nil {
		return nil, err
	}

	s.ctx.DispatchEvent(types.EventContainerCreated{})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return inst, nil
}

func (s *containerService) CheckForUpdates(ctx context.Context) (map[types.ContainerID]*types.Container, error) {
	for _, inst := range s.GetAll(ctx) {
		err := s.runnerService.CheckForUpdates(ctx, inst)
		if err != nil {
			return s.GetAll(ctx), err
		}
	}

	return s.GetAll(ctx), nil
}

func (s *containerService) load(ctx context.Context, uuid types.ContainerID) error {
	service, err := s.containerServiceService.Load(uuid)
	if err != nil {
		return err
	}

	inst := types.NewContainer(uuid, service)

	err = s.settingsService.Load(&inst)
	if err != nil {
		return err
	}

	err = s.envService.Load(&inst)
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

	if !s.Exists(ctx, uuid) {
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

func (s *containerService) SetDatabases(ctx context.Context, inst *types.Container, databases map[string]types.ContainerID, options map[string]*types.SetDatabasesOptions) error {
	for db := range databases {
		if _, ok := inst.Service.Databases[db]; !ok {
			return types.ErrDatabaseIDNotFound
		}
	}

	inst.Databases = databases
	err := s.settingsService.Save(inst, inst.ContainerSettings)
	if err != nil {
		return err
	}
	return s.remapDatabaseEnv(ctx, inst, options)
}

// remapDatabaseEnv remaps the environment variables of a container.
func (s *containerService) remapDatabaseEnv(ctx context.Context, inst *types.Container, options map[string]*types.SetDatabasesOptions) error {
	for databaseID, databaseContainerUUID := range inst.Databases {
		db, err := s.Get(ctx, databaseContainerUUID)
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

	return s.envService.Save(inst, inst.Env)
}

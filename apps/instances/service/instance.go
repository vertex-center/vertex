package service

import (
	"errors"
	"os"
	"path"
	"sync"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/adapter"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types/app"
)

var (
	ErrInstanceAlreadyExists      = errors.New("instance already exists")
	ErrInstanceAlreadyRunning     = errors.New("the instance is already running")
	ErrInstanceNotRunning         = errors.New("the instance is not running")
	ErrInstallMethodDoesNotExists = errors.New("this install method doesn't exist for this service")
)

type InstanceService struct {
	uuid uuid.UUID
	ctx  *app.Context

	instanceAdapter types.InstanceAdapterPort

	instanceRunnerService   *InstanceRunnerService
	instanceServiceService  *InstanceServiceService
	instanceEnvService      *InstanceEnvService
	instanceSettingsService *InstanceSettingsService

	instances      map[uuid.UUID]*types.Instance
	instancesMutex *sync.RWMutex
}

type InstanceServiceParams struct {
	Ctx *app.Context

	InstanceAdapter types.InstanceAdapterPort

	InstanceRunnerService   *InstanceRunnerService
	InstanceServiceService  *InstanceServiceService
	InstanceEnvService      *InstanceEnvService
	InstanceSettingsService *InstanceSettingsService
}

func NewInstanceService(params InstanceServiceParams) *InstanceService {
	s := &InstanceService{
		uuid: uuid.New(),
		ctx:  params.Ctx,

		instanceAdapter: params.InstanceAdapter,

		instanceRunnerService:   params.InstanceRunnerService,
		instanceServiceService:  params.InstanceServiceService,
		instanceEnvService:      params.InstanceEnvService,
		instanceSettingsService: params.InstanceSettingsService,

		instances:      make(map[uuid.UUID]*types.Instance),
		instancesMutex: &sync.RWMutex{},
	}

	s.ctx.AddListener(s)

	return s
}

// Get returns an instance by its UUID. If the instance does not exist,
// it returns ErrInstanceNotFound.
func (s *InstanceService) Get(uuid uuid.UUID) (*types.Instance, error) {
	s.instancesMutex.RLock()
	defer s.instancesMutex.RUnlock()

	instance, ok := s.instances[uuid]
	if !ok {
		return nil, types.ErrInstanceNotFound
	}
	return instance, nil
}

func (s *InstanceService) GetAll() map[uuid.UUID]*types.Instance {
	s.instancesMutex.RLock()
	defer s.instancesMutex.RUnlock()

	return s.instances
}

// Search returns all instances that match the query.
func (s *InstanceService) Search(query types.InstanceQuery) map[uuid.UUID]*types.Instance {
	instances := map[uuid.UUID]*types.Instance{}

	s.instancesMutex.RLock()
	defer s.instancesMutex.RUnlock()

	for _, inst := range s.instances {
		if !inst.HasOneOfFeatures(query.Features) {
			continue
		}

		instances[inst.UUID] = inst
	}

	return instances
}

func (s *InstanceService) Exists(uuid uuid.UUID) bool {
	s.instancesMutex.RLock()
	defer s.instancesMutex.RUnlock()

	return s.instances[uuid] != nil
}

// Delete deletes an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is still running, it returns ErrInstanceStillRunning.
func (s *InstanceService) Delete(uuid uuid.UUID) error {
	instance, err := s.Get(uuid)
	if err != nil {
		return err
	}

	serviceID := instance.Service.ID

	if instance.IsRunning() {
		return types.ErrInstanceStillRunning
	}

	err = s.instanceRunnerService.Delete(instance)
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	err = s.instanceAdapter.Delete(uuid)
	if err != nil {
		return err
	}

	s.instancesMutex.Lock()
	defer s.instancesMutex.Unlock()
	delete(s.instances, uuid)

	s.ctx.DispatchEvent(types.EventInstanceDeleted{
		InstanceUUID: uuid,
		ServiceID:    serviceID,
	})
	s.ctx.DispatchEvent(types.EventInstancesChange{})

	return nil
}

func (s *InstanceService) StartAll() {
	s.instancesMutex.RLock()
	defer s.instancesMutex.RUnlock()

	var ids []uuid.UUID

	for _, inst := range s.instances {
		if inst.LaunchOnStartup() {
			ids = append(ids, inst.UUID)
		}
	}

	if len(ids) == 0 {
		return
	}

	log.Info("trying to ping Google...")

	// Wait for internet connection
	if !wait.New(
		wait.WithWait(time.Second),
		wait.WithBreak(500*time.Millisecond),
	).Do([]string{"google.com:80"}) {
		log.Error(errors.New("internet connection: Failed to ping google.com"))
		return
	} else {
		log.Info("internet connection: OK")
	}

	// Start them
	for _, id := range ids {
		go func(id uuid.UUID) {
			inst, err := s.Get(id)
			if err != nil {
				log.Error(err)
				return
			}

			err = s.instanceRunnerService.Start(inst)
			if err != nil {
				log.Error(err)
			}
		}(id)
	}
}

func (s *InstanceService) StopAll() {
	s.instancesMutex.RLock()
	defer s.instancesMutex.RUnlock()

	for _, inst := range s.instances {
		err := s.instanceRunnerService.Stop(inst)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *InstanceService) Install(service types.Service, method string) (*types.Instance, error) {
	id := uuid.New()
	dir := s.instanceAdapter.GetPath(id)

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = s.instanceRunnerService.Install(id, service)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(path.Join(dir, ".vertex"), os.ModePerm)
	if err != nil {
		return nil, err
	}

	tempInstance := &types.Instance{
		UUID:    id,
		Service: service,
	}

	err = s.instanceServiceService.Save(tempInstance, service)
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

	inst.InstanceSettings.InstallMethod = &method
	err = s.instanceSettingsService.Save(inst, inst.InstanceSettings)
	if err != nil {
		return nil, err
	}

	inst.ResetDefaultEnv()
	err = s.instanceEnvService.Save(inst, inst.Env)
	if err != nil {
		return nil, err
	}

	s.ctx.DispatchEvent(types.EventInstanceCreated{})
	s.ctx.DispatchEvent(types.EventInstancesChange{})

	return inst, nil
}

func (s *InstanceService) CheckForUpdates() (map[uuid.UUID]*types.Instance, error) {
	for _, inst := range s.GetAll() {
		err := s.instanceRunnerService.CheckForUpdates(inst)
		if err != nil {
			return s.GetAll(), err
		}
	}

	return s.GetAll(), nil
}

func (s *InstanceService) LoadAll() {
	uuids, err := s.instanceAdapter.GetAll()
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

	s.ctx.DispatchEvent(types.EventInstancesLoaded{
		Count: loaded,
	})
}

func (s *InstanceService) load(uuid uuid.UUID) error {
	service, err := s.instanceServiceService.Load(uuid)
	if err != nil {
		return err
	}

	inst := types.NewInstance(uuid, service)

	err = s.instanceSettingsService.Load(&inst)
	if err != nil {
		return err
	}

	err = s.instanceEnvService.Load(&inst)
	if err != nil {
		return err
	}

	err = s.instanceServiceService.CheckForUpdate(&inst, service)
	if err != nil {
		return err
	}

	if !s.Exists(uuid) {
		s.instancesMutex.Lock()
		defer s.instancesMutex.Unlock()
		s.instances[uuid] = &inst
	} else {
		return ErrInstanceAlreadyExists
	}

	s.ctx.DispatchEvent(types.EventInstanceLoaded{
		Instance: &inst,
	})

	return nil
}

func (s *InstanceService) SetDatabases(inst *types.Instance, databases map[string]uuid.UUID) error {
	inst.Databases = databases
	err := s.instanceSettingsService.Save(inst, inst.InstanceSettings)
	if err != nil {
		return err
	}
	return s.remapDatabaseEnv(inst)
}

// remapDatabaseEnv remaps the environment variables of an instance.
func (s *InstanceService) remapDatabaseEnv(inst *types.Instance) error {
	for databaseID, databaseInstanceUUID := range inst.Databases {
		db, err := s.Get(databaseInstanceUUID)
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

	return s.instanceEnvService.Save(inst, inst.Env)
}

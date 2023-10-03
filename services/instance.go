package services

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

var (
	ErrInstanceAlreadyRunning     = errors.New("the instance is already running")
	ErrInstanceNotRunning         = errors.New("the instance is not running")
	ErrInstallMethodDoesNotExists = errors.New("this install method doesn't exist for this service")
)

type InstanceService struct {
	instanceAdapter     types.InstanceAdapterPort
	eventsAdapter       types.EventAdapterPort
	dockerRunnerAdapter types.RunnerAdapterPort
}

func NewInstanceService(instanceAdapter types.InstanceAdapterPort, dockerRunnerAdapter types.RunnerAdapterPort, eventRepo types.EventAdapterPort) InstanceService {
	s := InstanceService{
		instanceAdapter:     instanceAdapter,
		eventsAdapter:       eventRepo,
		dockerRunnerAdapter: dockerRunnerAdapter,
	}

	return s
}

// Get returns an instance by its UUID. If the instance does not exist,
// it returns ErrInstanceNotFound.
func (s *InstanceService) Get(uuid uuid.UUID) (*types.Instance, error) {
	return s.instanceAdapter.Get(uuid)
}

func (s *InstanceService) GetAll() map[uuid.UUID]*types.Instance {
	return s.instanceAdapter.GetAll()
}

// Search returns all instances that match the query.
func (s *InstanceService) Search(query types.InstanceQuery) map[uuid.UUID]*types.Instance {
	return s.instanceAdapter.Search(query)
}

// Delete deletes an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is still running, it returns ErrInstanceStillRunning.
func (s *InstanceService) Delete(uuid uuid.UUID) error {
	instance, err := s.instanceAdapter.Get(uuid)
	if err != nil {
		return err
	}

	if instance.IsRunning() {
		return types.ErrInstanceStillRunning
	}

	if instance.IsDockerized() {
		err = s.dockerRunnerAdapter.Delete(instance)
	} else {
		return fmt.Errorf("instance is not dockerized")
	}
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	err = s.instanceAdapter.Delete(uuid)
	if err != nil {
		return err
	}

	s.eventsAdapter.Send(types.EventInstancesChange{})
	return nil
}

// Start starts an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is already running, it returns ErrInstanceAlreadyRunning.
func (s *InstanceService) Start(uuid uuid.UUID) error {
	instance, err := s.instanceAdapter.Get(uuid)
	if err != nil {
		return err
	}

	if instance.IsBusy() {
		return nil
	}

	s.eventsAdapter.Send(types.EventInstanceLog{
		InstanceUUID: uuid,
		Kind:         types.LogKindOut,
		Message:      "Starting instance...",
	})

	log.Info("starting instance",
		vlog.String("uuid", uuid.String()),
	)

	if instance.IsRunning() {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexErr,
			Message:      ErrInstanceAlreadyRunning.Error(),
		})
		return ErrInstanceAlreadyRunning
	}

	setStatus := func(status string) {
		s.setStatus(instance, status)
	}

	var runner types.RunnerAdapterPort
	if instance.IsDockerized() {
		runner = s.dockerRunnerAdapter
	} else {
		return fmt.Errorf("instance is not dockerized")
	}

	stdout, stderr, err := runner.Start(instance, setStatus)
	if err != nil {
		s.setStatus(instance, types.InstanceStatusError)
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
					InstanceUUID: uuid,
					Kind:         types.LogKindDownload,
					Message:      strings.TrimPrefix(scanner.Text(), "DOWNLOAD"),
				})
				continue
			}

			s.eventsAdapter.Send(types.EventInstanceLog{
				InstanceUUID: uuid,
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
				InstanceUUID: uuid,
				Kind:         types.LogKindErr,
				Message:      scanner.Text(),
			})
		}
	}()

	// Wait for the instance until stopped
	wg.Wait()

	// Log stopped
	s.eventsAdapter.Send(types.EventInstanceLog{
		InstanceUUID: uuid,
		Kind:         types.LogKindVertexOut,
		Message:      "Stopping instance...",
	})
	log.Info("stopping instance",
		vlog.String("uuid", uuid.String()),
	)

	return nil
}

func (s *InstanceService) StartAll() {
	var ids []uuid.UUID

	for _, i := range s.instanceAdapter.GetAll() {
		if i.LaunchOnStartup() {
			ids = append(ids, i.UUID)
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
			err := s.Start(id)
			if err != nil {
				log.Error(err)
			}
		}(id)
	}
}

// Stop stops an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is not running, it returns ErrInstanceNotRunning.
func (s *InstanceService) Stop(uuid uuid.UUID) error {
	instance, err := s.instanceAdapter.Get(uuid)
	if err != nil {
		return err
	}

	if instance.IsBusy() {
		return nil
	}

	if !instance.IsRunning() {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexErr,
			Message:      ErrInstanceNotRunning.Error(),
		})
		return ErrInstanceNotRunning
	}

	s.setStatus(instance, types.InstanceStatusStopping)

	if instance.IsDockerized() {
		err = s.dockerRunnerAdapter.Stop(instance)
	} else {
		return fmt.Errorf("instance is not dockerized")
	}

	if err == nil {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexOut,
			Message:      "Instance stopped.",
		})

		log.Info("instance stopped",
			vlog.String("uuid", uuid.String()),
		)

		s.setStatus(instance, types.InstanceStatusOff)
	} else {
		s.setStatus(instance, types.InstanceStatusRunning)
	}

	return err
}

func (s *InstanceService) StopAll() {
	for _, i := range s.instanceAdapter.GetAll() {
		if !i.IsRunning() {
			continue
		}
		err := s.Stop(i.UUID)
		if err != nil {
			log.Error(err)
		}
	}
}

// WriteEnv writes environment variables to an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
func (s *InstanceService) WriteEnv(uuid uuid.UUID, environment map[string]string) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	i.Env = environment
	return s.instanceAdapter.SaveEnv(i)
}

func (s *InstanceService) Install(service types.Service, method string) (*types.Instance, error) {
	id := uuid.New()
	dir := s.instanceAdapter.GetPath(id)

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	if method == types.InstanceInstallMethodDocker {
		err = s.PreInstallForDocker(service, dir)
	} else {
		err = ErrInstallMethodDoesNotExists
	}
	if err != nil {
		return nil, err
	}

	tempInstance := &types.Instance{
		UUID:    id,
		Service: service,
	}

	err = s.instanceAdapter.SaveService(tempInstance)
	if err != nil {
		return nil, err
	}

	err = s.load(id)
	if err != nil {
		return nil, err
	}

	instance, err := s.instanceAdapter.Get(id)
	if err != nil {
		return nil, err
	}

	instance.InstanceSettings.InstallMethod = &method

	err = s.instanceAdapter.SaveSettings(instance)
	if err != nil {
		return nil, err
	}

	instance.ResetDefaultEnv()

	err = s.instanceAdapter.SaveEnv(instance)
	return instance, err
}

func (s *InstanceService) CheckForUpdates() (map[uuid.UUID]*types.Instance, error) {
	for _, instance := range s.GetAll() {
		var err error

		if instance.IsDockerized() {
			err = s.dockerRunnerAdapter.CheckForUpdates(instance)
		} else {
			return s.GetAll(), fmt.Errorf("instance is not dockerized")
		}

		if err != nil {
			return s.GetAll(), err
		}
	}

	return s.GetAll(), nil
}

// CheckForServiceUpdate checks if the service of an instance has an update.
// If it has, it sets the instance's ServiceUpdate field.
func (s *InstanceService) CheckForServiceUpdate(uuid uuid.UUID, latest types.Service) error {
	instance, err := s.Get(uuid)
	if err != nil {
		return err
	}

	current := instance.Service

	upToDate := reflect.DeepEqual(latest, current)
	log.Debug("service up-to-date", vlog.Bool("up_to_date", upToDate))
	instance.ServiceUpdate.Available = !upToDate
	return nil
}

// UpdateService updates the service of an instance by its UUID.
// The service passed is the latest version of the service.
func (s *InstanceService) UpdateService(uuid uuid.UUID, service types.Service) error {
	instance, err := s.Get(uuid)
	if err != nil {
		return err
	}

	if service.Version <= types.MaxSupportedVersion {
		log.Info("service version is outdated, upgrading.",
			vlog.String("uuid", uuid.String()),
			vlog.Int("old_version", int(instance.Service.Version)),
			vlog.Int("new_version", int(service.Version)),
		)
		instance.Service = service
		err := s.instanceAdapter.SaveService(instance)
		if err != nil {
			return err
		}

		err = s.CheckForServiceUpdate(uuid, service)
		if err != nil {
			return err
		}
	} else {
		log.Info("service version is not supported, skipping.",
			vlog.String("uuid", uuid.String()),
			vlog.Int("version", int(service.Version)),
		)
	}

	return nil
}

func (s *InstanceService) SetLaunchOnStartup(uuid uuid.UUID, value bool) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	i.InstanceSettings.LaunchOnStartup = &value
	return s.instanceAdapter.SaveSettings(i)
}

func (s *InstanceService) SetDisplayName(uuid uuid.UUID, value string) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	i.InstanceSettings.DisplayName = &value
	return s.instanceAdapter.SaveSettings(i)
}

func (s *InstanceService) SetDatabases(id uuid.UUID, databases map[string]uuid.UUID) error {
	i, err := s.Get(id)
	if err != nil {
		return err
	}

	i.Databases = databases

	err = s.remapDatabaseEnv(id)
	if err != nil {
		return err
	}

	return s.instanceAdapter.SaveSettings(i)
}

func (s *InstanceService) SetVersion(id uuid.UUID, value string) error {
	i, err := s.Get(id)
	if err != nil {
		return err
	}

	i.Version = &value
	return s.instanceAdapter.SaveSettings(i)
}

// remapDatabaseEnv remaps the environment variables of an instance.
func (s *InstanceService) remapDatabaseEnv(uuid uuid.UUID) error {
	instance, err := s.Get(uuid)
	if err != nil {
		return err
	}

	for databaseID, databaseInstanceUUID := range instance.Databases {
		db, err := s.Get(databaseInstanceUUID)
		if err != nil {
			return err
		}

		host := config.Current.HostVertex
		if strings.Contains(host, ":") {
			host, _, err = net.SplitHostPort(host)
			if err != nil {
				return err
			}
		}

		dbEnvNames := (*db.Service.Features.Databases)[0]
		iEnvNames := instance.Service.Databases[databaseID].Names

		instance.Env[iEnvNames.Host] = host
		instance.Env[iEnvNames.Port] = db.Env[dbEnvNames.Port]
		if dbEnvNames.Username != nil {
			instance.Env[iEnvNames.Username] = db.Env[*dbEnvNames.Username]
		}
		if dbEnvNames.Password != nil {
			instance.Env[iEnvNames.Password] = db.Env[*dbEnvNames.Password]
		}
	}

	err = s.instanceAdapter.SaveEnv(instance)
	if err != nil {
		return err
	}

	return s.instanceAdapter.SaveSettings(instance)
}

func (s *InstanceService) GetDockerContainerInfo(uuid uuid.UUID) (map[string]any, error) {
	instance, err := s.Get(uuid)
	if err != nil {
		return nil, err
	}

	if !instance.IsDockerized() {
		return nil, errors.New("instance is not using docker")
	}

	info, err := s.dockerRunnerAdapter.Info(*instance)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *InstanceService) GetAllVersions(instance *types.Instance, useCache bool) ([]string, error) {
	if !useCache || len(instance.CacheVersions) == 0 {
		versions, err := s.dockerRunnerAdapter.GetAllVersions(*instance)
		if err != nil {
			return nil, err
		}
		instance.CacheVersions = versions
	}

	return instance.CacheVersions, nil
}

func (s *InstanceService) CloneRepository(dir string, url string) error {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

func (s *InstanceService) DownloadRelease(dir string, repo string) error {
	split := strings.Split(repo, "/")

	owner := split[1]
	repository := split[2]

	return storage.DownloadLatestGithubRelease(owner, repository, dir)
}

func (s *InstanceService) Symlink(path string, repo string) error {
	p := strings.Split(repo, ":")[1]

	_, err := s.instanceAdapter.LoadService(p)
	if err != nil {
		return fmt.Errorf("%s is not a compatible Vertex service", repo)
	}

	return os.Symlink(p, path)
}

func (s *InstanceService) setStatus(instance *types.Instance, status string) {
	if instance.Status == status {
		return
	}

	instance.Status = status
	s.eventsAdapter.Send(types.EventInstancesChange{})
	s.eventsAdapter.Send(types.EventInstanceStatusChange{
		InstanceUUID: instance.UUID,
		Status:       status,
	})
}

func (s *InstanceService) LoadAll() {
	s.instanceAdapter.Reload(func(uuid uuid.UUID) {
		err := s.load(uuid)
		if err != nil {
			log.Error(err)
			return
		}
	})
}

func (s *InstanceService) load(uuid uuid.UUID) error {
	instancePath := s.instanceAdapter.GetPath(uuid)

	service, err := s.instanceAdapter.LoadService(instancePath)
	if err != nil {
		return err
	}

	instance := types.NewInstance(uuid, service)

	err = s.instanceAdapter.LoadSettings(&instance)
	if err != nil {
		return err
	}

	err = s.instanceAdapter.LoadEnv(&instance)
	if err != nil {
		return err
	}

	err = s.instanceAdapter.Set(uuid, instance)
	if err != nil {
		return err
	}

	err = s.CheckForServiceUpdate(uuid, service)
	if err != nil {
		return err
	}

	s.eventsAdapter.Send(types.EventInstanceLoaded{
		InstanceUuid: uuid,
	})

	return nil
}

func (s *InstanceService) PreInstallForDocker(service types.Service, dir string) error {
	if service.Methods.Docker == nil {
		return ErrInstallMethodDoesNotExists
	}

	if service.Methods.Docker.Clone != nil {
		err := s.CloneRepository(dir, service.Methods.Docker.Clone.Repository)
		if err != nil {
			return err
		}
	}

	return nil
}

// RecreateContainer recreates a container by its UUID.
func (s *InstanceService) RecreateContainer(instance *types.Instance) error {
	if !instance.IsDockerized() {
		return nil
	}

	if instance.IsRunning() {
		err := s.dockerRunnerAdapter.Stop(instance)
		if err != nil {
			return err
		}
	}

	err := s.dockerRunnerAdapter.Delete(instance)
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	go func() {
		err := s.Start(instance.UUID)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	return nil
}

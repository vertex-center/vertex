package services

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v2"
)

var (
	ErrInstanceAlreadyRunning     = errors.New("the instance is already running")
	ErrInstanceNotRunning         = errors.New("the instance is not running")
	ErrInstallMethodDoesNotExists = errors.New("this install method doesn't exist for this service")
)

type InstanceService struct {
	uuid uuid.UUID

	serviceAdapter  types.ServiceAdapterPort
	instanceAdapter types.InstanceAdapterPort
	logsAdapter     types.InstanceLogsAdapterPort
	eventsAdapter   types.EventAdapterPort

	dockerRunnerAdapter types.RunnerAdapterPort
	fsRunnerAdapter     types.RunnerAdapterPort
}

func NewInstanceService(serviceAdapter types.ServiceAdapterPort, instanceAdapter types.InstanceAdapterPort, dockerRunnerAdapter types.RunnerAdapterPort, fsRunnerAdapter types.RunnerAdapterPort, instanceLogsAdapter types.InstanceLogsAdapterPort, eventRepo types.EventAdapterPort) InstanceService {
	s := InstanceService{
		uuid: uuid.New(),

		serviceAdapter:      serviceAdapter,
		instanceAdapter:     instanceAdapter,
		logsAdapter:         instanceLogsAdapter,
		eventsAdapter:       eventRepo,
		dockerRunnerAdapter: dockerRunnerAdapter,
		fsRunnerAdapter:     fsRunnerAdapter,
	}

	s.reload()

	eventRepo.AddListener(&s)

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
		err = s.fsRunnerAdapter.Delete(instance)
	}
	if err != nil {
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

	onLog := func(msg string) {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindOut,
			Message:      msg,
		})
	}

	onErr := func(msg string) {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindErr,
			Message:      msg,
		})
	}

	setStatus := func(status string) {
		s.setStatus(instance, status)
	}

	if instance.IsDockerized() {
		err = s.dockerRunnerAdapter.Start(instance, onLog, onErr, setStatus)
	} else {
		err = s.fsRunnerAdapter.Start(instance, onLog, onErr, setStatus)
	}

	if err != nil {
		s.setStatus(instance, types.InstanceStatusError)
	} else {
		s.eventsAdapter.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexOut,
			Message:      "Instance started.",
		})

		log.Info("instance started",
			vlog.String("uuid", uuid.String()),
		)
	}

	return err
}

func (s *InstanceService) StartAll() {
	var ids []uuid.UUID

	// Select instances that should launch on startup
	for _, i := range s.instanceAdapter.GetAll() {
		launchOnStartup := i.InstanceSettings.LaunchOnStartup

		if launchOnStartup != nil && !*launchOnStartup {
			continue
		}

		ids = append(ids, i.UUID)
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

	s.eventsAdapter.Send(types.EventInstanceLog{
		InstanceUUID: uuid,
		Kind:         types.LogKindVertexOut,
		Message:      "Stopping instance...",
	})

	log.Info("stopping instance",
		vlog.String("uuid", uuid.String()),
	)

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
		err = s.fsRunnerAdapter.Stop(instance)
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

	return s.instanceAdapter.SaveEnv(i, environment)
}

func (s *InstanceService) Install(serviceID string, method string) (*types.Instance, error) {
	id := uuid.New()
	dir := s.instanceAdapter.GetPath(id)

	service, err := s.serviceAdapter.Get(serviceID)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	if method == types.InstanceInstallMethodScript {
		err = s.PreInstallForScript(service, dir)
	} else if method == types.InstanceInstallMethodRelease {
		err = s.PreInstallForRelease(service, dir)
	} else if method == types.InstanceInstallMethodDocker {
		err = s.PreInstallForDocker(service, dir)
	}
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(path.Join(dir, ".vertex"), os.ModePerm)
	if err != nil {
		return nil, err
	}

	serviceYaml, err := yaml.Marshal(service)
	if err != nil {

		return nil, err
	}

	err = os.WriteFile(path.Join(dir, ".vertex", "service.yml"), serviceYaml, os.ModePerm)
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

	env := map[string]string{}
	for _, v := range instance.Service.Env {
		env[v.Name] = v.Default
	}

	err = s.instanceAdapter.SaveEnv(instance, env)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *InstanceService) CheckForUpdates() (map[uuid.UUID]*types.Instance, error) {
	for _, instance := range s.GetAll() {
		var err error

		if instance.IsDockerized() {
			err = s.dockerRunnerAdapter.CheckForUpdates(instance)
		} else {
			err = s.fsRunnerAdapter.CheckForUpdates(instance)
		}

		if err != nil {
			return s.GetAll(), err
		}
	}

	return s.GetAll(), nil
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

// remapDatabaseEnv remaps the environment variables of an instance.
func (s *InstanceService) remapDatabaseEnv(uuid uuid.UUID) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	for databaseID, databaseInstanceUUID := range i.Databases {
		db, err := s.Get(databaseInstanceUUID)
		if err != nil {
			return err
		}

		host := config.Current.Host
		if strings.Contains(host, ":") {
			host, _, err = net.SplitHostPort(host)
			if err != nil {
				return err
			}
		}

		dbEnvNames := (*db.Service.Features.Databases)[0]
		iEnvNames := i.Service.Databases[databaseID].Names

		i.Env[iEnvNames.Host] = host
		i.Env[iEnvNames.Port] = db.Env[dbEnvNames.Port]
		if dbEnvNames.Username != nil {
			i.Env[iEnvNames.Username] = db.Env[*dbEnvNames.Username]
		}
		if dbEnvNames.Password != nil {
			i.Env[iEnvNames.Password] = db.Env[*dbEnvNames.Password]
		}
	}

	err = s.instanceAdapter.SaveEnv(i, i.Env)
	if err != nil {
		return err
	}

	return s.instanceAdapter.SaveSettings(i)
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

func (s *InstanceService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.logsAdapter.LoadBuffer(uuid)
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

	_, err := s.instanceAdapter.ReadService(p)
	if err != nil {
		return fmt.Errorf("%s is not a compatible Vertex service", repo)
	}

	return os.Symlink(p, path)
}

func (s *InstanceService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceLog:
		s.logsAdapter.Push(e.InstanceUUID, types.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	}
}

func (s *InstanceService) GetUUID() uuid.UUID {
	return s.uuid
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

func (s *InstanceService) reload() {
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

	service, err := s.instanceAdapter.ReadService(instancePath)
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

	err = s.logsAdapter.Open(uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *InstanceService) PreInstallForScript(service types.Service, dir string) error {
	if service.Methods.Script == nil {
		return ErrInstallMethodDoesNotExists
	}

	if service.Methods.Script.Clone != nil {
		err := s.CloneRepository(dir, service.Methods.Script.Clone.Repository)
		if err != nil {
			return err
		}
	}

	script, err := s.serviceAdapter.GetScript(service.ID)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, service.Methods.Script.Filename), script, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *InstanceService) PreInstallForRelease(service types.Service, dir string) error {
	if service.Methods.Release == nil {
		return ErrInstallMethodDoesNotExists
	}

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
	if err != nil {
		return err
	}
	return s.Start(instance.UUID)
}

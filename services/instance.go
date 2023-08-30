package services

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"gopkg.in/yaml.v2"
)

var (
	ErrInstanceAlreadyRunning     = errors.New("the instance is already running")
	ErrInstanceNotRunning         = errors.New("the instance is not running")
	ErrInstallMethodDoesNotExists = errors.New("this install method doesn't exist for this service")
)

type InstanceService struct {
	uuid uuid.UUID

	serviceRepo  types.ServiceRepository
	instanceRepo types.InstanceRepository
	logsRepo     types.InstanceLogsRepository
	eventsRepo   types.EventRepository

	dockerRunnerRepo types.RunnerRepository
	fsRunnerRepo     types.RunnerRepository
}

func NewInstanceService(serviceRepo types.ServiceRepository, instanceRepo types.InstanceRepository, dockerRunnerRepo types.RunnerRepository, fsRunnerRepo types.RunnerRepository, instanceLogsRepo types.InstanceLogsRepository, eventRepo types.EventRepository) InstanceService {
	s := InstanceService{
		uuid: uuid.New(),

		serviceRepo:      serviceRepo,
		instanceRepo:     instanceRepo,
		logsRepo:         instanceLogsRepo,
		eventsRepo:       eventRepo,
		dockerRunnerRepo: dockerRunnerRepo,
		fsRunnerRepo:     fsRunnerRepo,
	}

	s.reload()

	eventRepo.AddListener(&s)

	return s
}

func (s *InstanceService) Get(uuid uuid.UUID) (*types.Instance, error) {
	return s.instanceRepo.Get(uuid)
}

func (s *InstanceService) GetAll() map[uuid.UUID]*types.Instance {
	return s.instanceRepo.GetAll()
}

func (s *InstanceService) Delete(uuid uuid.UUID) error {
	instance, err := s.instanceRepo.Get(uuid)
	if err != nil {
		return err
	}

	if instance.IsRunning() {
		return errors.New("failed to delete this instance because the instance is still running")
	}

	if instance.IsDockerized() {
		err = s.dockerRunnerRepo.Delete(instance)
	} else {
		err = s.fsRunnerRepo.Delete(instance)
	}
	if err != nil {
		return err
	}

	err = s.instanceRepo.Delete(uuid)
	if err != nil {
		return err
	}

	s.eventsRepo.Send(types.EventInstancesChange{})
	return nil
}

func (s *InstanceService) Start(uuid uuid.UUID) error {
	instance, err := s.instanceRepo.Get(uuid)
	if err != nil {
		return err
	}

	s.eventsRepo.Send(types.EventInstanceLog{
		InstanceUUID: uuid,
		Kind:         types.LogKindOut,
		Message:      "Starting instance...",
	})

	logger.Log("starting instance").
		AddKeyValue("uuid", uuid).
		Print()

	if instance.IsRunning() {
		s.eventsRepo.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexErr,
			Message:      ErrInstanceAlreadyRunning.Error(),
		})
		return ErrInstanceAlreadyRunning
	}

	onLog := func(msg string) {
		s.eventsRepo.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindOut,
			Message:      msg,
		})
	}

	onErr := func(msg string) {
		s.eventsRepo.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindErr,
			Message:      msg,
		})
	}

	setStatus := func(status string) {
		s.setStatus(instance, status)
	}

	if instance.IsDockerized() {
		err = s.dockerRunnerRepo.Start(instance, onLog, onErr, setStatus)
	} else {
		err = s.fsRunnerRepo.Start(instance, onLog, onErr, setStatus)
	}

	if err != nil {
		s.setStatus(instance, types.InstanceStatusError)
	} else {
		s.eventsRepo.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexOut,
			Message:      "Instance started.",
		})

		logger.Log("instance started").
			AddKeyValue("uuid", uuid).
			Print()
	}

	return err
}

func (s *InstanceService) StartAll() {
	var ids []uuid.UUID

	// Select instances that should launch on startup
	for _, i := range s.instanceRepo.GetAll() {
		launchOnStartup := i.InstanceSettings.LaunchOnStartup

		if launchOnStartup != nil && !*launchOnStartup {
			continue
		}

		ids = append(ids, i.UUID)
	}

	if len(ids) == 0 {
		return
	}

	logger.Log("trying to ping Google...").Print()

	// Wait for internet connection
	if !wait.New(
		wait.WithWait(time.Second),
		wait.WithBreak(500*time.Millisecond),
	).Do([]string{"google.com:80"}) {
		logger.Error(errors.New("internet connection: Failed to ping google.com")).Print()
		return
	} else {
		logger.Log("internet connection: OK").Print()
	}

	// Start them
	for _, id := range ids {
		go func(id uuid.UUID) {
			err := s.Start(id)
			if err != nil {
				logger.Error(err).Print()
			}
		}(id)
	}
}

func (s *InstanceService) Stop(uuid uuid.UUID) error {
	instance, err := s.instanceRepo.Get(uuid)
	if err != nil {
		return err
	}

	s.eventsRepo.Send(types.EventInstanceLog{
		InstanceUUID: uuid,
		Kind:         types.LogKindVertexOut,
		Message:      "Stopping instance...",
	})

	logger.Log("stopping instance").
		AddKeyValue("uuid", uuid).
		Print()

	if !instance.IsRunning() {
		s.eventsRepo.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexErr,
			Message:      ErrInstanceNotRunning.Error(),
		})
		return ErrInstanceNotRunning
	}

	s.setStatus(instance, types.InstanceStatusStopping)

	if instance.IsDockerized() {
		err = s.dockerRunnerRepo.Stop(instance)
	} else {
		err = s.fsRunnerRepo.Stop(instance)
	}

	if err == nil {
		s.eventsRepo.Send(types.EventInstanceLog{
			InstanceUUID: uuid,
			Kind:         types.LogKindVertexOut,
			Message:      "Instance stopped.",
		})

		logger.Log("instance stopped").
			AddKeyValue("uuid", uuid).
			Print()

		s.setStatus(instance, types.InstanceStatusOff)
	} else {
		s.setStatus(instance, types.InstanceStatusRunning)
	}

	return err
}

func (s *InstanceService) StopAll() {
	for _, i := range s.instanceRepo.GetAll() {
		if !i.IsRunning() {
			continue
		}
		err := s.Stop(i.UUID)
		if err != nil {
			logger.Error(err).Print()
		}
	}
}

func (s *InstanceService) WriteEnv(uuid uuid.UUID, environment map[string]string) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	return s.instanceRepo.SaveEnv(i, environment)
}

func (s *InstanceService) Install(serviceID string, method string) (*types.Instance, error) {
	id := uuid.New()
	dir := path.Join(storage.PathInstances, id.String())

	service, err := s.serviceRepo.Get(serviceID)
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

	instance, err := s.instanceRepo.Get(id)
	if err != nil {
		return nil, err
	}

	instance.InstanceSettings.InstallMethod = &method

	err = s.instanceRepo.SaveSettings(instance)
	if err != nil {
		return nil, err
	}

	env := map[string]string{}
	for _, v := range instance.EnvDefinitions {
		env[v.Name] = v.Default
	}

	err = s.instanceRepo.SaveEnv(instance, env)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *InstanceService) CheckForUpdates() (map[uuid.UUID]*types.Instance, error) {
	for _, instance := range s.GetAll() {
		var err error

		if instance.IsDockerized() {
			err = s.dockerRunnerRepo.CheckForUpdates(instance)
		} else {
			err = s.fsRunnerRepo.CheckForUpdates(instance)
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
	return s.instanceRepo.SaveSettings(i)
}

func (s *InstanceService) SetDisplayName(uuid uuid.UUID, value string) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	i.InstanceSettings.DisplayName = &value
	return s.instanceRepo.SaveSettings(i)
}

func (s *InstanceService) GetDockerContainerInfo(uuid uuid.UUID) (map[string]any, error) {
	instance, err := s.Get(uuid)
	if err != nil {
		return nil, err
	}

	if !instance.IsDockerized() {
		return nil, errors.New("instance is not using docker")
	}

	info, err := s.dockerRunnerRepo.Info(*instance)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *InstanceService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.logsRepo.LoadBuffer(uuid)
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

	_, err := s.instanceRepo.ReadService(p)
	if err != nil {
		return fmt.Errorf("%s is not a compatible Vertex service", repo)
	}

	return os.Symlink(p, path)
}

func (s *InstanceService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceLog:
		s.logsRepo.Push(e.InstanceUUID, types.LogLine{
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
	s.eventsRepo.Send(types.EventInstancesChange{})
	s.eventsRepo.Send(types.EventInstanceStatusChange{
		InstanceUUID: instance.UUID,
		Status:       status,
	})
}

func (s *InstanceService) reload() {
	s.instanceRepo.Reload(func(uuid uuid.UUID) {
		err := s.load(uuid)
		if err != nil {
			logger.Error(err).Print()
			return
		}
	})
}

func (s *InstanceService) load(uuid uuid.UUID) error {
	instancePath := path.Join(storage.PathInstances, uuid.String())

	service, err := s.instanceRepo.ReadService(instancePath)
	if err != nil {
		return err
	}

	instance := types.NewInstance(uuid, service)

	err = s.instanceRepo.LoadSettings(&instance)
	if err != nil {
		return err
	}

	err = s.instanceRepo.LoadEnv(&instance)
	if err != nil {
		return err
	}

	err = s.instanceRepo.Set(uuid, instance)
	if err != nil {
		return err
	}

	err = s.logsRepo.Open(uuid)
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

	script, err := s.serviceRepo.GetScript(service.ID)
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

func (s *InstanceService) RecreateContainer(instance *types.Instance) error {
	if !instance.IsDockerized() {
		return nil
	}

	if instance.IsRunning() {
		err := s.dockerRunnerRepo.Stop(instance)
		if err != nil {
			return err
		}
	}

	err := s.dockerRunnerRepo.Delete(instance)
	if err != nil {
		return err
	}
	return s.Start(instance.UUID)
}

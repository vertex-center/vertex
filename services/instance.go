package services

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/nakabonne/tstorage"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/repository"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrContainerStillRunning  = errors.New("the container is still running")
	ErrInstanceAlreadyRunning = errors.New("the instance is already running")
	ErrInstanceNotRunning     = errors.New("the instance is not running")
)

type InstanceService struct {
	repo       repository.InstanceRepository
	dockerRepo repository.DockerRepository
}

func NewInstanceService() InstanceService {
	return InstanceService{
		repo:       repository.NewInstanceRepository(),
		dockerRepo: repository.NewDockerRepository(),
	}
}

func (s *InstanceService) Unload() {
	s.repo.Unload()
}

func (s *InstanceService) Get(uuid uuid.UUID) (*types.Instance, error) {
	return s.repo.Get(uuid)
}

func (s *InstanceService) GetAll() map[uuid.UUID]*types.Instance {
	return s.repo.GetAll()
}

func (s *InstanceService) Delete(uuid uuid.UUID) error {
	i, err := s.repo.Get(uuid)
	if err != nil {
		return err
	}

	if i.IsRunning() {
		return ErrContainerStillRunning
	}

	if i.UseDocker {
		containerID, err := s.dockerRepo.GetContainerID(i.DockerContainerName())
		if err == repository.ErrContainerNotFound {
			logger.Warn(err.Error()).Print()
		} else if err != nil {
			return err
		} else {
			err = s.dockerRepo.RemoveContainer(containerID)
			if err != nil {
				return err
			}
		}
	}

	return s.repo.Delete(uuid)
}

func (s *InstanceService) AddListener(channel chan types.InstanceEvent) uuid.UUID {
	return s.repo.AddListener(channel)
}

func (s *InstanceService) RemoveListener(uuid uuid.UUID) {
	s.repo.RemoveListener(uuid)
}

func (s *InstanceService) Start(uuid uuid.UUID) error {
	i, err := s.repo.Get(uuid)
	if err != nil {
		return err
	}

	s.repo.WriteLogLine(i, &types.LogLine{
		Kind:    types.LogKindVertexOut,
		Message: "Starting instance...",
	})

	logger.Log("starting instance").
		AddKeyValue("uuid", uuid).
		Print()

	if i.IsRunning() {
		s.repo.WriteLogLine(i, &types.LogLine{
			Kind:    types.LogKindVertexErr,
			Message: ErrInstanceAlreadyRunning.Error(),
		})
		return ErrInstanceAlreadyRunning
	}

	if i.UseDocker {
		err = s.startWithDocker(i)
	} else {
		err = s.startManually(i)
	}

	if err != nil {
		i.SetStatus(types.InstanceStatusError)
	} else {
		s.repo.WriteLogLine(i, &types.LogLine{
			Kind:    types.LogKindVertexOut,
			Message: "Instance started.",
		})

		logger.Log("instance started").
			AddKeyValue("uuid", uuid).
			Print()

		s.startUptimeRoutine(i)
	}

	return err
}

func (s *InstanceService) startUptimeRoutine(i *types.Instance) {
	i.UptimeStopChannels = []*chan bool{}
	for _, url := range i.URLs {
		go func(name string, url string) {
			ch := make(chan bool, 1)
			i.UptimeStopChannels = append(i.UptimeStopChannels, &ch)
			ticker := time.NewTicker(time.Second * 5)

			defer func() {
				_ = i.PushStatus(name, types.UptimeStatusFloatOff)
				ticker.Stop()
				close(ch)
				logger.Log("uptime ticker stopped").
					AddKeyValue("instance_uuid", i.UUID).
					Print()
			}()

			for {
				select {
				case <-ch:
					return
				case <-ticker.C:
					client := http.Client{
						Timeout: time.Second * 2,
					}
					res, err := client.Get(url)
					if err != nil {
						logger.Error(err).Print()
						break
					}
					if res.StatusCode >= 400 {
						err = i.PushStatus(name, types.UptimeStatusFloatOff)
					} else {
						err = i.PushStatus(name, types.UptimeStatusFloatOn)
					}
					if err != nil {
						logger.Error(err).Print()
					}
					res.Body.Close()
				}
			}
		}(url.Name, "http://localhost:"+url.Port+*url.PingRoute)
	}
}

func (s *InstanceService) StartAll() {
	for _, i := range s.repo.GetAll() {
		launchOnStartup := i.InstanceMetadata.LaunchOnStartup
		if launchOnStartup != nil && !*launchOnStartup {
			continue
		}
		err := s.Start(i.UUID)
		if err != nil {
			logger.Error(err).Print()
		}
	}
}

func (s *InstanceService) Stop(uuid uuid.UUID) error {
	i, err := s.repo.Get(uuid)
	if err != nil {
		return err
	}

	s.repo.WriteLogLine(i, &types.LogLine{
		Kind:    types.LogKindVertexOut,
		Message: "Stopping instance...",
	})

	logger.Log("stopping instance").
		AddKeyValue("uuid", uuid).
		Print()

	if !i.IsRunning() {
		s.repo.WriteLogLine(i, &types.LogLine{
			Kind:    types.LogKindVertexErr,
			Message: ErrInstanceNotRunning.Error(),
		})
		return ErrInstanceNotRunning
	}

	s.stopUptimeRoutine(i)

	if i.UseDocker {
		err = s.stopWithDocker(i)
	} else {
		err = s.stopManually(i)
	}

	if err == nil {
		s.repo.WriteLogLine(i, &types.LogLine{
			Kind:    types.LogKindVertexOut,
			Message: "Instance stopped.",
		})

		logger.Log("instance stopped").
			AddKeyValue("uuid", uuid).
			Print()

		i.SetStatus(types.InstanceStatusOff)
	}

	return err
}

func (s *InstanceService) stopUptimeRoutine(i *types.Instance) {
	for _, ch := range i.UptimeStopChannels {
		*ch <- true
	}
}

func (s *InstanceService) StopAll() {
	for _, i := range s.repo.GetAll() {
		if !i.IsRunning() {
			continue
		}
		err := s.Stop(i.UUID)
		if err != nil {
			logger.Error(err).Print()
		}
	}
}

func (s *InstanceService) startWithDocker(i *types.Instance) error {
	imageName := i.DockerImageName()
	containerName := i.DockerContainerName()

	i.SetStatus(types.InstanceStatusBuilding)

	instancePath := s.repo.GetPath(i)

	onMsg := func(msg string) {
		s.repo.WriteLogLine(i, &types.LogLine{
			Kind:    types.LogKindOut,
			Message: msg,
		})
	}

	// Build
	var err error
	if i.Methods.Docker.Dockerfile != nil {
		err = s.dockerRepo.BuildImageFromDockerfile(instancePath, imageName, onMsg)
	} else if i.Methods.Docker.Image != nil {
		err = s.dockerRepo.BuildImageFromName(*i.Methods.Docker.Image, onMsg)
	} else {
		return errors.New("no Docker methods found")
	}

	if err != nil {
		s.repo.WriteLogLine(i, &types.LogLine{
			Kind:    types.LogKindErr,
			Message: err.Error(),
		})
		return err
	}

	// Create
	id, err := s.dockerRepo.GetContainerID(containerName)
	if err == repository.ErrContainerNotFound {
		logger.Log("container doesn't exists, create it.").
			AddKeyValue("container_name", containerName).
			Print()

		exposedPorts := nat.PortSet{}
		portBindings := nat.PortMap{}
		if i.Methods.Docker.Ports != nil {
			var all []string

			for _, out := range *i.Methods.Docker.Ports {
				in := ""
				for _, e := range i.EnvDefinitions {
					if e.Type == "port" && e.Default == out {
						in = i.EnvVariables.Entries[e.Name]
						all = append(all, in+":"+out)
						break
					}
				}
			}

			var err error
			exposedPorts, portBindings, err = nat.ParsePortSpecs(all)
			if err != nil {
				return err
			}
		}

		var binds []string
		if i.Methods.Docker.Volumes != nil {
			for source, target := range *i.Methods.Docker.Volumes {
				source, err = filepath.Abs(path.Join(instancePath, "volumes", source))
				if err != nil {
					return err
				}
				binds = append(binds, source+":"+target)
			}
		}

		if i.Methods.Docker.Dockerfile != nil {
			id, err = s.dockerRepo.CreateContainer(imageName, containerName, exposedPorts, portBindings, binds)
		} else if i.Methods.Docker.Image != nil {
			id, err = s.dockerRepo.CreateContainer(*i.Methods.Docker.Image, containerName, exposedPorts, portBindings, binds)
		}
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	i.SetStatus(types.InstanceStatusStarting)

	// Start
	err = s.dockerRepo.StartContainer(id)
	if err != nil {
		return err
	}

	i.SetStatus(types.InstanceStatusRunning)
	return nil
}

func (s *InstanceService) startManually(i *types.Instance) error {
	if i.Cmd != nil {
		logger.Error(errors.New("runner already started")).
			AddKeyValue("name", i.Name).
			Print()
	}

	dir := s.repo.GetPath(i)
	executable := i.ID
	command := "./" + i.ID

	// Try to find the executable
	// For a service of ID=vertex-id, the executable can be:
	// - vertex-id
	// - script-filename.sh
	_, err := os.Stat(path.Join(dir, executable))
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Stat(path.Join(dir, i.Methods.Script.Filename))
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("executables %s and %s were not found", i.ID, i.Methods.Script.Filename)
		} else if err != nil {
			return err
		}
		command = fmt.Sprintf("./%s", i.Methods.Script.Filename)
	} else if err != nil {
		return err
	}

	i.Cmd = exec.Command(command)
	i.Cmd.Dir = dir

	i.Cmd.Stdin = os.Stdin

	stdoutReader, err := i.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := i.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			s.repo.WriteLogLine(i, &types.LogLine{
				Kind:    types.LogKindOut,
				Message: stdoutScanner.Text(),
			})
		}
	}()

	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			s.repo.WriteLogLine(i, &types.LogLine{
				Kind:    types.LogKindErr,
				Message: stderrScanner.Text(),
			})
		}
	}()

	i.SetStatus(types.InstanceStatusRunning)

	err = i.Cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := i.Cmd.Wait()
		if err != nil {
			logger.Error(err).
				AddKeyValue("name", i.Service.Name).
				Print()
		}
		i.SetStatus(types.InstanceStatusOff)
	}()

	return nil
}

func (s *InstanceService) stopWithDocker(i *types.Instance) error {
	id, err := s.dockerRepo.GetContainerID(i.DockerContainerName())
	if err != nil {
		return err
	}

	return s.dockerRepo.StopContainer(id)
}

func (s *InstanceService) stopManually(i *types.Instance) error {
	err := i.Cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Force kill if the process continues

	i.Cmd = nil

	return nil
}

func (s *InstanceService) WriteEnv(uuid uuid.UUID, environment map[string]string) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	return s.repo.WriteEnv(i, environment)
}

func (s *InstanceService) Install(repo string, useDocker *bool, useReleases *bool) (*types.Instance, error) {
	serviceUUID := uuid.New()
	basePath := path.Join(storage.PathInstances, serviceUUID.String())

	forceClone := (useDocker != nil && *useDocker) || (useReleases == nil || !*useReleases)

	var err error
	if strings.HasPrefix(repo, "marketplace:") {
		err = s.repo.Download(basePath, repo, forceClone)
	} else if strings.HasPrefix(repo, "localstorage:") {
		err = s.repo.Symlink(basePath, repo)
	} else if strings.HasPrefix(repo, "git:") {
		err = s.repo.Download(basePath, repo, forceClone)
	} else {
		return nil, fmt.Errorf("this protocol is not supported")
	}

	if err != nil {
		return nil, err
	}

	i, err := s.repo.Instantiate(serviceUUID)
	if err != nil {
		return nil, err
	}

	if useDocker != nil {
		i.InstanceMetadata.UseDocker = *useDocker
	}
	if useReleases != nil {
		i.InstanceMetadata.UseReleases = *useReleases
	}

	err = s.repo.SaveMetadata(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *InstanceService) SetLaunchOnStartup(uuid uuid.UUID, value bool) error {
	i, err := s.Get(uuid)
	if err != nil {
		return err
	}

	i.InstanceMetadata.LaunchOnStartup = &value
	err = s.repo.SaveMetadata(i)
	if err != nil {
		return err
	}

	return err
}

func (s *InstanceService) GetDockerContainerInfo(uuid uuid.UUID) (*types.DockerContainerInfo, error) {
	i, err := s.Get(uuid)
	if err != nil {
		return nil, err
	}

	if !i.UseDocker {
		return nil, errors.New("instance is not using docker")
	}

	info, err := s.dockerRepo.GetContainerInfo(i.DockerContainerName())
	if err != nil {
		return nil, err
	}

	return info, nil
}

type StatusSince int

const (
	StatusSinceOneHour = iota
	StatusSinceTwoDay
)

func (s *InstanceService) GetAllStatus(uuid uuid.UUID, since StatusSince) ([]types.Uptime, error) {
	i, err := s.Get(uuid)
	if err != nil {
		return nil, err
	}

	var (
		uptimes  []types.Uptime
		from     time.Time
		count    int
		interval time.Duration
	)

	switch since {
	case StatusSinceTwoDay:
		from = time.Now().Add(-time.Hour * 48)
		count = 48
		interval = time.Hour
	case StatusSinceOneHour:
		from = time.Now().Add(-time.Hour)
		count = 60
		interval = time.Minute
	}

	for _, url := range i.URLs {
		if url.PingRoute == nil {
			continue
		}

		var (
			history                 []types.UptimePoint
			currentStatusFloat      float64 = -1
			currentRangeStatusFloat float64
		)

		t := from
		for j := 0; j < count; j += 1 {
			currentRangeStatusFloat = currentStatusFloat

			start := t
			end := start.Add(interval)

			points, err := i.UptimeStorage.Select(
				"status_change",
				[]tstorage.Label{{Name: "name", Value: url.Name}},
				start.Unix(),
				end.Unix(),
			)
			if err != nil && err != tstorage.ErrNoDataPoints {
				return nil, err
			}

			for _, p := range points {
				currentStatusFloat = p.Value
				if currentRangeStatusFloat > p.Value {
					currentRangeStatusFloat = p.Value
				}
			}

			history = append(history, types.UptimePoint{
				Status: types.UptimeStatus(currentRangeStatusFloat),
			})

			t = end
		}

		uptimes = append(uptimes, types.Uptime{
			Name:    url.Name,
			PingURL: url.PingRoute,
			Current: types.UptimeStatus(currentStatusFloat),
			History: history,
		})
	}

	return uptimes, nil
}

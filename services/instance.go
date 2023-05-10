package services

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/nakabonne/tstorage"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrContainerStillRunning  = errors.New("the container is still running")
	ErrInstanceAlreadyRunning = errors.New("the instance is already running")
	ErrInstanceNotRunning     = errors.New("the instance is not running")
)

type InstanceService struct {
	instanceRepo types.InstanceRepository

	dockerRunnerRepo types.RunnerRepository
	fsRunnerRepo     types.RunnerRepository
}

func NewInstanceService(instanceRepo types.InstanceRepository, dockerRunnerRepo types.RunnerRepository, fsRunnerRepo types.RunnerRepository) InstanceService {
	return InstanceService{
		instanceRepo:     instanceRepo,
		dockerRunnerRepo: dockerRunnerRepo,
		fsRunnerRepo:     fsRunnerRepo,
	}
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
		return ErrContainerStillRunning
	}

	if instance.IsDockerized() {
		err = s.dockerRunnerRepo.Delete(instance)
	} else {
		err = s.fsRunnerRepo.Delete(instance)
	}
	if err != nil {
		return err
	}

	return s.instanceRepo.Delete(uuid)
}

func (s *InstanceService) AddListener(channel chan types.InstanceEvent) uuid.UUID {
	return s.instanceRepo.AddListener(channel)
}

func (s *InstanceService) RemoveListener(uuid uuid.UUID) {
	s.instanceRepo.RemoveListener(uuid)
}

func (s *InstanceService) Start(uuid uuid.UUID) error {
	instance, err := s.instanceRepo.Get(uuid)
	if err != nil {
		return err
	}

	instance.PushLogLine(&types.LogLine{
		Kind:    types.LogKindVertexOut,
		Message: "Starting instance...",
	})

	logger.Log("starting instance").
		AddKeyValue("uuid", uuid).
		Print()

	if instance.IsRunning() {
		instance.PushLogLine(&types.LogLine{
			Kind:    types.LogKindVertexErr,
			Message: ErrInstanceAlreadyRunning.Error(),
		})
		return ErrInstanceAlreadyRunning
	}

	if instance.IsDockerized() {
		err = s.dockerRunnerRepo.Start(instance)
	} else {
		err = s.fsRunnerRepo.Start(instance)
	}

	if err != nil {
		instance.SetStatus(types.InstanceStatusError)
	} else {
		instance.PushLogLine(&types.LogLine{
			Kind:    types.LogKindVertexOut,
			Message: "Instance started.",
		})

		logger.Log("instance started").
			AddKeyValue("uuid", uuid).
			Print()

		s.startUptimeRoutine(instance)
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
	for _, i := range s.instanceRepo.GetAll() {
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
	instance, err := s.instanceRepo.Get(uuid)
	if err != nil {
		return err
	}

	instance.PushLogLine(&types.LogLine{
		Kind:    types.LogKindVertexOut,
		Message: "Stopping instance...",
	})

	logger.Log("stopping instance").
		AddKeyValue("uuid", uuid).
		Print()

	if !instance.IsRunning() {
		instance.PushLogLine(&types.LogLine{
			Kind:    types.LogKindVertexErr,
			Message: ErrInstanceNotRunning.Error(),
		})
		return ErrInstanceNotRunning
	}

	s.stopUptimeRoutine(instance)

	if instance.IsDockerized() {
		err = s.dockerRunnerRepo.Stop(instance)
	} else {
		err = s.fsRunnerRepo.Stop(instance)
	}

	if err == nil {
		instance.PushLogLine(&types.LogLine{
			Kind:    types.LogKindVertexOut,
			Message: "Instance stopped.",
		})

		logger.Log("instance stopped").
			AddKeyValue("uuid", uuid).
			Print()

		instance.SetStatus(types.InstanceStatusOff)
	}

	return err
}

func (s *InstanceService) stopUptimeRoutine(i *types.Instance) {
	for _, ch := range i.UptimeStopChannels {
		*ch <- true
	}
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

func (s *InstanceService) Install(repo string, useDocker *bool, useReleases *bool) (*types.Instance, error) {
	serviceUUID := uuid.New()
	basePath := path.Join(storage.PathInstances, serviceUUID.String())

	forceClone := (useDocker != nil && *useDocker) || (useReleases == nil || !*useReleases)

	var err error
	if strings.HasPrefix(repo, "marketplace:") {
		err = s.Download(basePath, repo, forceClone)
	} else if strings.HasPrefix(repo, "localstorage:") {
		err = s.Symlink(basePath, repo)
	} else if strings.HasPrefix(repo, "git:") {
		err = s.Download(basePath, repo, forceClone)
	} else {
		return nil, fmt.Errorf("this protocol is not supported")
	}

	if err != nil {
		return nil, err
	}

	i, err := s.instanceRepo.Load(serviceUUID)
	if err != nil {
		return nil, err
	}

	i.InstanceMetadata.UseDocker = useDocker
	i.InstanceMetadata.UseReleases = useReleases

	err = s.instanceRepo.SaveMetadata(i)
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
	err = s.instanceRepo.SaveMetadata(i)
	if err != nil {
		return err
	}

	return err
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
		uptimes   []types.Uptime
		from      time.Time
		count     int
		interval  time.Duration
		remaining int
	)

	switch since {
	case StatusSinceTwoDay:
		from = time.Now().Add(-time.Hour * 48).Truncate(time.Hour)
		count = 48
		interval = time.Hour
		remaining = 3600 - time.Now().Hour()
	case StatusSinceOneHour:
		from = time.Now().Add(-time.Hour).Truncate(time.Minute)
		count = 60
		interval = time.Minute
		remaining = 60 - time.Now().Second()
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
			Name:             url.Name,
			PingURL:          url.PingRoute,
			Current:          types.UptimeStatus(currentStatusFloat),
			IntervalSeconds:  int(interval.Seconds()),
			RemainingSeconds: remaining,
			History:          history,
		})
	}

	return uptimes, nil
}

func (s *InstanceService) Download(dest string, repo string, forceClone bool) error {
	var err error

	if forceClone {
		logger.Log("force-clone enabled.").Print()
	} else {
		logger.Log("force-clone disabled. try to download the releases first").Print()
		err = downloadFromReleases(dest, repo)
	}

	if forceClone || errors.Is(err, storage.ErrNoReleasesPublished) {
		split := strings.Split(repo, ":")
		repo = "git:https://" + split[1]

		err = downloadFromGit(dest, repo)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *InstanceService) Symlink(path string, repo string) error {
	p := strings.Split(repo, ":")[1]

	_, err := s.instanceRepo.ReadService(p)
	if err != nil {
		return fmt.Errorf("%s is not a compatible Vertex service", repo)
	}

	return os.Symlink(p, path)
}

func downloadFromReleases(dest string, repo string) error {
	split := strings.Split(repo, "/")

	owner := split[1]
	repository := split[2]

	return storage.DownloadLatestGithubRelease(owner, repository, dest)
}

func downloadFromGit(path string, repo string) error {
	url := strings.SplitN(repo, ":", 2)[1]
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

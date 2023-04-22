package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	logger = console.New("vertex::repository")

	ErrContainerNotFound = errors.New("container not found")
)

const (
	EventStdout = "stdout"
	EventStderr = "stderr"
	EventChange = "change"
)

type InstanceRepository struct {
	instances map[uuid.UUID]*types.Instance
	listeners map[uuid.UUID]chan types.InstanceEvent
	observer  chan types.InstanceEvent
}

func NewInstanceRepository() InstanceRepository {
	r := InstanceRepository{
		instances: map[uuid.UUID]*types.Instance{},
		listeners: map[uuid.UUID]chan types.InstanceEvent{},
		observer:  make(chan types.InstanceEvent),
	}

	r.reload()

	go func() {
		defer close(r.observer)

		for {
			<-r.observer
			r.notifyListeners(types.InstanceEvent{
				Name: EventChange,
			})
		}
	}()

	return r
}

func (r *InstanceRepository) GetPath(i *types.Instance) string {
	return path.Join(storage.PathInstances, i.UUID.String())
}

func (r *InstanceRepository) Get(uuid uuid.UUID) (*types.Instance, error) {
	i := r.instances[uuid]
	if i == nil {
		return nil, fmt.Errorf("the service '%s' is not instances", uuid)
	}
	return i, nil
}

func (r *InstanceRepository) GetAll() map[uuid.UUID]*types.Instance {
	return r.instances
}

func (r *InstanceRepository) Delete(uuid uuid.UUID) error {
	i := r.instances[uuid]

	err := os.RemoveAll(r.GetPath(i))
	if err != nil {
		return fmt.Errorf("failed to delete server uuid=%s: %v", i.UUID, err)
	}

	delete(r.instances, uuid)

	r.notifyListeners(types.InstanceEvent{
		Name: EventChange,
	})

	return nil
}

func (r *InstanceRepository) Exists(uuid uuid.UUID) bool {
	return r.instances[uuid] != nil
}

func (r *InstanceRepository) Create(uuid uuid.UUID, i *types.Instance) {
	r.instances[uuid] = i
	r.notifyListeners(types.InstanceEvent{
		Name: EventChange,
	})
}

func (r *InstanceRepository) AddListener(channel chan types.InstanceEvent) uuid.UUID {
	id := uuid.New()
	r.listeners[id] = channel
	logger.Log(fmt.Sprintf("channel %s registered to instances", id))
	return id
}

func (r *InstanceRepository) RemoveListener(uuid uuid.UUID) {
	delete(r.listeners, uuid)
	logger.Log(fmt.Sprintf("channel %s unregistered from instances", uuid))
}

func (r *InstanceRepository) SaveMetadata(i *types.Instance) error {
	metaPath := path.Join(r.GetPath(i), ".vertex", "instance_metadata.json")

	metaBytes, err := json.MarshalIndent(i.InstanceMetadata, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(metaPath, metaBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (r *InstanceRepository) notifyListeners(event types.InstanceEvent) {
	for _, listener := range r.listeners {
		listener <- event
	}
}

func (r *InstanceRepository) Instantiate(uuid uuid.UUID) (*types.Instance, error) {
	if r.Exists(uuid) {
		return nil, fmt.Errorf("the service '%s' is already running", uuid)
	}

	i, err := r.load(uuid)
	if err != nil {
		return nil, err
	}

	r.Create(uuid, i)

	i.Register(r.observer)

	return i, nil
}

func (r *InstanceRepository) reload() {
	entries, err := os.ReadDir(storage.PathInstances)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Fatal(err)
		}

		isInstance := entry.IsDir() || info.Mode()&os.ModeSymlink != 0

		if isInstance {
			logger.Log(fmt.Sprintf("found service uuid=%s", entry.Name()))
			serviceUUID, err := uuid.Parse(entry.Name())
			if err != nil {
				log.Fatal(err)
			}

			if !r.Exists(serviceUUID) {
				logger.Log(fmt.Sprintf("instantiate service uuid=%s", entry.Name()))

				_, err = r.Instantiate(serviceUUID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func (r *InstanceRepository) load(instanceUUID uuid.UUID) (*types.Instance, error) {
	service, err := r.readService(path.Join(storage.PathInstances, instanceUUID.String()))
	if err != nil {
		return nil, err
	}

	meta := types.InstanceMetadata{
		UseDocker:   false,
		UseReleases: false,
	}

	metaPath := path.Join(storage.PathInstances, instanceUUID.String(), ".vertex", "instance_metadata.json")
	metaBytes, err := os.ReadFile(metaPath)

	if errors.Is(err, os.ErrNotExist) {
		logger.Log("instance_metadata.json not found. using default.")
	} else if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(metaBytes, &meta)
		if err != nil {
			return nil, err
		}
	}

	i := &types.Instance{
		Service:          *service,
		InstanceMetadata: meta,
		Status:           types.InstanceStatusOff,
		Logs:             types.Logs{},
		EnvVariables:     *types.NewEnvVariables(),
		UUID:             instanceUUID,
		Listeners:        map[uuid.UUID]chan types.InstanceEvent{},
	}

	err = r.readEnv(i)
	return i, err
}

func (r *InstanceRepository) readService(servicePath string) (*types.Service, error) {
	data, err := os.ReadFile(path.Join(servicePath, ".vertex", "service.json"))
	if err != nil {
		logger.Warn(fmt.Sprintf("service at '%s' has no '.vertex/service.json' file", path.Dir(servicePath)))
	}

	var service types.Service
	err = json.Unmarshal(data, &service)
	return &service, err
}

func (r *InstanceRepository) readEnv(i *types.Instance) error {
	filepath := path.Join(r.GetPath(i), ".env")

	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) < 2 {
			return errors.New("failed to read .env")
		}

		i.EnvVariables.Entries[line[0]] = line[1]
	}

	return nil
}

func (r *InstanceRepository) WriteEnv(i *types.Instance, variables map[string]string) error {
	filepath := path.Join(r.GetPath(i), ".env")

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	for key, value := range variables {
		_, err := file.WriteString(strings.Join([]string{key, value}, "=") + "\n")
		if err != nil {
			return err
		}
	}

	i.EnvVariables.Entries = variables

	return nil
}

func (r *InstanceRepository) Symlink(path string, repo string) error {
	p := strings.Split(repo, ":")[1]

	_, err := r.readService(p)
	if err != nil {
		return fmt.Errorf("%s is not a compatible Vertex service", repo)
	}

	return os.Symlink(p, path)
}

func (r *InstanceRepository) Download(dest string, repo string, forceClone bool) error {
	var err error

	if forceClone {
		logger.Log("force-clone enabled.")
	} else {
		logger.Log("force-clone disabled. try to download the releases first")
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

func (r *InstanceRepository) AppendLogLine(i *types.Instance, line *types.LogLine) {
	i.Logs.Add(line)

	data, err := json.Marshal(line)
	if err != nil {
		logger.Error(err)
	}

	var name string
	switch line.Kind {
	case EventStderr:
		name = EventStderr
	default:
		name = EventStdout
	}

	i.NotifyListeners(types.InstanceEvent{
		Name: name,
		Data: string(data),
	})
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

package instances

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/services/instance"
	"github.com/vertex-center/vertex/storage"
)

var (
	instancesObserver chan instance.Event
	logger            = console.New("vertex::services-manager")
)

var instances = newInstances()

const (
	EventChange = "change"
)

type Event struct {
	Name string
}

type Instances struct {
	all map[uuid.UUID]*instance.Instance

	listeners map[uuid.UUID]chan Event
}

func newInstances() *Instances {
	instancesObserver = make(chan instance.Event)

	instances := &Instances{
		all:       map[uuid.UUID]*instance.Instance{},
		listeners: map[uuid.UUID]chan Event{},
	}

	go func() {
		defer close(instancesObserver)

		for {
			select {
			case _ = <-instancesObserver:
				for _, listener := range instances.listeners {
					listener <- Event{
						Name: EventChange,
					}
				}
			}
		}
	}()

	return instances
}

func Start(uuid uuid.UUID) error {
	i, err := Get(uuid)
	if err != nil {
		return err
	}
	return i.Start()
}

func Stop(uuid uuid.UUID) error {
	i, err := Get(uuid)
	if err != nil {
		return err
	}
	return i.Stop()
}

func Add(uuid uuid.UUID, i *instance.Instance) {
	instances.all[uuid] = i
	for _, listener := range instances.listeners {
		listener <- Event{
			Name: EventChange,
		}
	}
}

func Delete(uuid uuid.UUID) error {
	i, err := Get(uuid)
	if err != nil {
		return err
	}

	err = i.Delete()
	if err != nil {
		return err
	}

	delete(instances.all, uuid)

	for _, listener := range instances.listeners {
		listener <- Event{
			Name: EventChange,
		}
	}

	return nil
}

func Exists(uuid uuid.UUID) bool {
	return instances.all[uuid] != nil
}

func Instantiate(uuid uuid.UUID) (*instance.Instance, error) {
	if Exists(uuid) {
		return nil, fmt.Errorf("the service '%s' is already running", uuid)
	}

	i, err := instance.CreateFromDisk(uuid)
	if err != nil {
		return nil, err
	}

	Add(uuid, i)

	i.Register(instancesObserver)

	return i, nil
}

func List() map[uuid.UUID]*instance.Instance {
	return instances.all
}

func Get(uuid uuid.UUID) (*instance.Instance, error) {
	i := instances.all[uuid]
	if i == nil {
		return nil, fmt.Errorf("the service '%s' is not instances", uuid)
	}
	return i, nil
}

func Register(channel chan Event) uuid.UUID {
	id := uuid.New()
	instances.listeners[id] = channel
	logger.Log(fmt.Sprintf("channel %s registered to instances", id))
	return id
}

func Unregister(uuid uuid.UUID) {
	delete(instances.listeners, uuid)
	logger.Log(fmt.Sprintf("channel %s unregistered from instances", uuid))
}

func Install(repo string) (*instance.Instance, error) {
	serviceUUID := uuid.New()
	basePath := path.Join(storage.PathInstances, serviceUUID.String())

	var err error
	if strings.HasPrefix(repo, "marketplace:") {
		err = downloadFromMarketplace(basePath, repo)
		// If there are no releases available, that may mean that the repository
		// should be cloned. This happens when there are repositories that don't
		// need to compile things, but only run a bash script.
		if errors.Is(err, storage.ErrNoReleasesPublished) {
			split := strings.Split(repo, ":")
			repo = "git:https://" + split[1]

			err = downloadFromGit(basePath, repo)
			if err != nil {
				return nil, err
			}
		}
	} else if strings.HasPrefix(repo, "localstorage:") {
		err = downloadFromLocalstorage(basePath, repo)
	} else if strings.HasPrefix(repo, "git:") {
		err = downloadFromGit(basePath, repo)
	} else {
		return nil, fmt.Errorf("this protocol is not supported")
	}

	if err != nil {
		return nil, err
	}

	return Instantiate(serviceUUID)
}

func downloadFromMarketplace(dest string, repo string) error {
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

func downloadFromLocalstorage(path string, repo string) error {
	p := strings.Split(repo, ":")[1]

	_, err := services.ReadFromDisk(p)
	if err != nil {
		return fmt.Errorf("%s is not a compatible Vertex service", repo)
	}

	return os.Symlink(p, path)
}

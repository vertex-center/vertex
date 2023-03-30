package instances

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
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

	if strings.HasPrefix(repo, "marketplace:") {
		client := github.NewClient(nil)

		split := strings.Split(repo, "/")

		owner := split[1]
		repo := split[2]

		release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to retrieve the latest github release for %s: %v", repo, err))
		}

		platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

		for _, asset := range release.Assets {
			if strings.Contains(*asset.Name, platform) {
				archivePath := path.Join(basePath, "temp.tar.gz")

				err := downloadFile(*asset.BrowserDownloadURL, basePath, archivePath)
				if err != nil {
					return nil, err
				}

				err = untarFile(basePath, archivePath)
				if err != nil {
					return nil, err
				}

				err = os.Remove(archivePath)
				if err != nil {
					return nil, err
				}

				break
			}
		}
	} else if strings.HasPrefix(repo, "localstorage:") {
		p := strings.Split(repo, ":")[1]

		_, err := services.ReadFromDisk(p)
		if err != nil {
			return nil, fmt.Errorf("%s is not a compatible Vertex service", basePath)
		}

		err = os.Symlink(p, basePath)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(repo, "git:") {
		url := strings.SplitN(repo, ":", 2)[1]

		_, err := git.PlainClone(basePath, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("this protocol is not supported")
	}

	i, err := Instantiate(serviceUUID)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func downloadFile(url string, basePath string, archivePath string) error {
	err := os.Mkdir(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func untarFile(basePath string, archivePath string) error {
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	stream, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer stream.Close()

	reader := tar.NewReader(stream)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		filepath := path.Join(basePath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			err = os.Chmod(filepath, 0755)
			if err != nil {
				return err
			}

			file.Close()
		default:
			return errors.New(fmt.Sprintf("unknown flag type (%s) for file '%s'", header.Typeflag, header.Name))
		}
	}

	return nil
}

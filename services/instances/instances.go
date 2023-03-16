package instances

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services"
)

var logger = console.New("vertex::services-manager")

var instances = Instances{
	all:       map[uuid.UUID]*Instance{},
	listeners: map[uuid.UUID]chan Event{},
}

const (
	StatusOff     = "off"
	StatusRunning = "running"
	StatusError   = "error"
)

const (
	EventChange = "change"
)

type Event struct {
	Name string
}

type Instances struct {
	all map[uuid.UUID]*Instance

	listeners map[uuid.UUID]chan Event
}

type Instance struct {
	services.Service

	Status string `json:"status"`

	uuid uuid.UUID
	cmd  *exec.Cmd
}

func Start(uuid uuid.UUID) error {
	instance, err := Get(uuid)
	if err != nil {
		return err
	}

	if instance.cmd != nil {
		logger.Error(fmt.Errorf("runner %s already started", instance.Name))
	}

	instance.cmd = exec.Command(fmt.Sprintf("./%s", instance.ID))
	instance.cmd.Dir = path.Join("servers", instance.uuid.String())

	instance.cmd.Stdout = os.Stdout
	instance.cmd.Stderr = os.Stderr
	instance.cmd.Stdin = os.Stdin

	setStatus(instance, StatusRunning)

	err = instance.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := instance.cmd.Wait()
		if err != nil {
			logger.Error(fmt.Errorf("%s: %v", instance.Service.Name, err))
		}
		setStatus(instance, StatusOff)
	}()

	return nil
}

func Stop(uuid uuid.UUID) error {
	instance, err := Get(uuid)
	if err != nil {
		return err
	}

	err = instance.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Remove runner from runners
	// TODO: Force kill if the process continues

	instance.cmd = nil
	return nil
}

func CreateFromDisk(uuid uuid.UUID) (*Instance, error) {
	data, err := os.ReadFile(path.Join("servers", uuid.String(), ".vertex", "service.json"))
	if err != nil {
		logger.Warn(fmt.Sprintf("service '%s' has no '.vertex/service.json' file", uuid))
	}

	var service services.Service
	err = json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Service: service,
		Status:  StatusOff,
		uuid:    uuid,
	}, nil
}

func Add(uuid uuid.UUID, instance *Instance) {
	instances.all[uuid] = instance
	for _, listener := range instances.listeners {
		listener <- Event{
			Name: EventChange,
		}
	}
}

func Exists(uuid uuid.UUID) bool {
	return instances.all[uuid] != nil
}

func Instantiate(uuid uuid.UUID) (*Instance, error) {
	if Exists(uuid) {
		return nil, fmt.Errorf("the service '%s' is already running", uuid)
	}

	instance, err := CreateFromDisk(uuid)
	if err != nil {
		return nil, err
	}

	Add(uuid, instance)

	return instance, nil
}

func List() map[uuid.UUID]*Instance {
	return instances.all
}

func Get(uuid uuid.UUID) (*Instance, error) {
	instance := instances.all[uuid]
	if instance == nil {
		return nil, fmt.Errorf("the service '%s' is not instances", uuid)
	}
	return instance, nil
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

func Install(s services.Service) (*Instance, error) {
	if strings.HasPrefix(s.Repository, "github") {
		client := github.NewClient(nil)

		split := strings.Split(s.Repository, "/")

		owner := split[1]
		repo := split[2]

		release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to retrieve the latest github release for %s", s.Repository))
		}

		platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

		serviceUUID := uuid.New()

		for _, asset := range release.Assets {
			if strings.Contains(*asset.Name, platform) {
				basePath := path.Join("servers", serviceUUID.String())
				archivePath := path.Join(basePath, fmt.Sprintf("%s.tar.gz", s.ID))

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

		instance, err := Instantiate(serviceUUID)
		if err != nil {
			return nil, err
		}

		return instance, nil
	}

	return nil, errors.New("this repository is not supported")
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

func setStatus(instance *Instance, status string) {
	instance.Status = status
	for _, listener := range instances.listeners {
		listener <- Event{
			Name: EventChange,
		}
	}
}

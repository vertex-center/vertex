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

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrContainerNotFound = errors.New("container not found")
)

const (
	EventChange = "change"
)

type InstanceFSRepository struct {
	instances map[uuid.UUID]*types.Instance
}

func NewInstanceFSRepository() InstanceFSRepository {
	r := InstanceFSRepository{
		instances: map[uuid.UUID]*types.Instance{},
	}

	return r
}

func (r *InstanceFSRepository) Get(uuid uuid.UUID) (*types.Instance, error) {
	i := r.instances[uuid]
	if i == nil {
		return nil, fmt.Errorf("the service '%s' is not instances", uuid)
	}
	return i, nil
}

func (r *InstanceFSRepository) GetAll() map[uuid.UUID]*types.Instance {
	return r.instances
}

func (r *InstanceFSRepository) GetPath(uuid uuid.UUID) string {
	return path.Join(storage.PathInstances, uuid.String())
}

func (r *InstanceFSRepository) Delete(uuid uuid.UUID) error {
	err := os.RemoveAll(r.GetPath(uuid))
	if err != nil {
		return errors.New("failed to delete server")
	}

	delete(r.instances, uuid)

	return nil
}

func (r *InstanceFSRepository) Exists(uuid uuid.UUID) bool {
	return r.instances[uuid] != nil
}

func (r *InstanceFSRepository) Set(uuid uuid.UUID, instance types.Instance) error {
	if r.Exists(uuid) {
		return fmt.Errorf("the instance '%s' already exists", uuid)
	}

	r.instances[uuid] = &instance

	return nil
}

func (r *InstanceFSRepository) SaveMetadata(i *types.Instance) error {
	metaPath := path.Join(r.GetPath(i.UUID), ".vertex", "instance_metadata.json")

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

func (r *InstanceFSRepository) LoadMetadata(i *types.Instance) error {
	metaPath := path.Join(r.GetPath(i.UUID), ".vertex", "instance_metadata.json")
	metaBytes, err := os.ReadFile(metaPath)

	if errors.Is(err, os.ErrNotExist) {
		logger.Log("instance_metadata.json not found. using default.").Print()
	} else if err != nil {
		return err
	} else {
		err = json.Unmarshal(metaBytes, &i.InstanceMetadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InstanceFSRepository) ReadService(instancePath string) (types.Service, error) {
	data, err := os.ReadFile(path.Join(instancePath, ".vertex", "service.json"))
	if err != nil {
		logger.Warn("service has no '.vertex/service.json' file").
			AddKeyValue("path", path.Dir(instancePath)).
			Print()
	}

	var service types.Service
	err = json.Unmarshal(data, &service)
	return service, err
}

func (r *InstanceFSRepository) SaveEnv(i *types.Instance, variables map[string]string) error {
	filepath := path.Join(r.GetPath(i.UUID), ".env")

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

	i.EnvVariables = variables

	return nil
}

func (r *InstanceFSRepository) LoadEnv(i *types.Instance) error {
	filepath := path.Join(r.GetPath(i.UUID), ".env")

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

		i.EnvVariables[line[0]] = line[1]
	}

	return nil
}

func (r *InstanceFSRepository) Close() {
	for _, instance := range r.instances {
		err := instance.UptimeStorage.Close()
		if err != nil {
			logger.Error(err).Print()
		}
	}
}

func (r *InstanceFSRepository) Reload(load func(uuid uuid.UUID)) {
	r.Close()

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
			logger.Log("found service").
				AddKeyValue("uuid", entry.Name()).
				Print()

			id, err := uuid.Parse(entry.Name())
			if err != nil {
				log.Fatal(err)
			}

			load(id)
		}
	}
}

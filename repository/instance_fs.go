package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
	"gopkg.in/yaml.v2"
)

var (
	ErrInstanceNotFound      = errors.New("instance not found")
	ErrInstanceAlreadyExists = errors.New("instance already exists")
	ErrContainerNotFound     = errors.New("container not found")
)

type InstanceFSRepository struct {
	instancesPath string
	instances     map[uuid.UUID]*types.Instance
}

func NewInstanceFSRepository() InstanceFSRepository {
	r := InstanceFSRepository{
		instancesPath: path.Join(storage.Path, "instances"),
		instances:     map[uuid.UUID]*types.Instance{},
	}

	err := os.MkdirAll(r.instancesPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Default.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", r.instancesPath),
		)
		os.Exit(1)
	}

	return r
}

func (r *InstanceFSRepository) Get(uuid uuid.UUID) (*types.Instance, error) {
	instance, ok := r.instances[uuid]
	if !ok {
		return nil, ErrInstanceNotFound
	}
	return instance, nil
}

func (r *InstanceFSRepository) GetAll() map[uuid.UUID]*types.Instance {
	return r.instances
}

func (r *InstanceFSRepository) GetPath(uuid uuid.UUID) string {
	return path.Join(r.instancesPath, uuid.String())
}

func (r *InstanceFSRepository) Delete(uuid uuid.UUID) error {
	err := os.RemoveAll(r.GetPath(uuid))
	if err != nil {
		return fmt.Errorf("failed to delete server: %v", err)
	}

	delete(r.instances, uuid)

	return nil
}

func (r *InstanceFSRepository) Exists(uuid uuid.UUID) bool {
	return r.instances[uuid] != nil
}

func (r *InstanceFSRepository) Set(uuid uuid.UUID, instance types.Instance) error {
	if r.Exists(uuid) {
		return ErrInstanceAlreadyExists
	}

	r.instances[uuid] = &instance

	return nil
}

func (r *InstanceFSRepository) SaveSettings(i *types.Instance) error {
	settingsPath := path.Join(r.GetPath(i.UUID), ".vertex", "instance_settings.json")

	settingsBytes, err := json.MarshalIndent(i.InstanceSettings, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(settingsPath, settingsBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (r *InstanceFSRepository) LoadSettings(i *types.Instance) error {
	settingsPath := path.Join(r.GetPath(i.UUID), ".vertex", "instance_settings.json")
	settingsBytes, err := os.ReadFile(settingsPath)

	if errors.Is(err, os.ErrNotExist) {
		log.Default.Warn("instance_settings.json not found. using default.")
	} else if err != nil {
		return err
	} else {
		err = json.Unmarshal(settingsBytes, &i.InstanceSettings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InstanceFSRepository) ReadService(instancePath string) (types.Service, error) {
	data, err := os.ReadFile(path.Join(instancePath, ".vertex", "service.yml"))
	if err != nil {
		log.Default.Warn("'.vertex/service.yml' file not found",
			vlog.String("path", path.Dir(instancePath)),
		)
	}

	var service types.Service
	err = yaml.Unmarshal(data, &service)
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

func (r *InstanceFSRepository) Reload(load func(uuid uuid.UUID)) {
	entries, err := os.ReadDir(r.instancesPath)
	if err != nil {
		log.Default.Error(err)
		os.Exit(1)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Default.Error(err)
			continue
		}

		isInstance := entry.IsDir() || info.Mode()&os.ModeSymlink != 0

		if isInstance {
			log.Default.Info("found instance",
				vlog.String("uuid", entry.Name()),
			)

			id, err := uuid.Parse(entry.Name())
			if err != nil {
				log.Default.Error(err)
				continue
			}

			load(id)
		}
	}
}

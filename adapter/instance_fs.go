package adapter

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
	"gopkg.in/yaml.v3"
)

type InstanceFilePath string

const (
	InstanceVertexPath   InstanceFilePath = ".vertex"
	InstanceSettingsPath InstanceFilePath = ".vertex/instance_settings.json"
	InstanceServicePath  InstanceFilePath = ".vertex/service.yml"
	InstanceEnvPath      InstanceFilePath = ".env"
)

var (
	ErrInstanceAlreadyExists = errors.New("instance already exists")
	ErrContainerNotFound     = errors.New("container not found")
)

type InstanceFSAdapter struct {
	instancesPath string
	instances     map[uuid.UUID]*types.Instance
}

type InstanceFSAdapterParams struct {
	instancesPath string
}

func NewInstanceFSAdapter(params *InstanceFSAdapterParams) types.InstanceAdapterPort {
	if params == nil {
		params = &InstanceFSAdapterParams{}
	}
	if params.instancesPath == "" {
		params.instancesPath = path.Join(storage.Path, "instances")
	}

	adapter := &InstanceFSAdapter{
		instancesPath: params.instancesPath,
		instances:     map[uuid.UUID]*types.Instance{},
	}

	err := os.MkdirAll(adapter.instancesPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", adapter.instancesPath),
		)
		os.Exit(1)
	}

	return adapter
}

// Get returns an instance by its UUID. If the instance does not exist,
// it returns ErrInstanceNotFound.
func (a *InstanceFSAdapter) Get(uuid uuid.UUID) (*types.Instance, error) {
	instance, ok := a.instances[uuid]
	if !ok {
		return nil, types.ErrInstanceNotFound
	}
	return instance, nil
}

func (a *InstanceFSAdapter) GetAll() map[uuid.UUID]*types.Instance {
	return a.instances
}

func (a *InstanceFSAdapter) Search(query types.InstanceQuery) map[uuid.UUID]*types.Instance {
	instances := map[uuid.UUID]*types.Instance{}

	for _, instance := range a.instances {
		if !instance.HasOneOfFeatures(query.Features) {
			continue
		}

		instances[instance.UUID] = instance
	}

	return instances
}

func (a *InstanceFSAdapter) GetPath(uuid uuid.UUID) string {
	return path.Join(a.instancesPath, uuid.String())
}

func (a *InstanceFSAdapter) GetFilePath(uuid uuid.UUID, filepath InstanceFilePath) string {
	return path.Join(a.GetPath(uuid), string(filepath))
}

func (a *InstanceFSAdapter) Delete(uuid uuid.UUID) error {
	err := os.RemoveAll(a.GetPath(uuid))
	if err != nil {
		return fmt.Errorf("failed to delete server: %v", err)
	}

	delete(a.instances, uuid)

	return nil
}

func (a *InstanceFSAdapter) Exists(uuid uuid.UUID) bool {
	return a.instances[uuid] != nil
}

func (a *InstanceFSAdapter) Set(uuid uuid.UUID, instance types.Instance) error {
	if a.Exists(uuid) {
		return ErrInstanceAlreadyExists
	}

	a.instances[uuid] = &instance

	return nil
}

func (a *InstanceFSAdapter) SaveSettings(i *types.Instance) error {
	settingsPath := a.GetFilePath(i.UUID, InstanceSettingsPath)

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

func (a *InstanceFSAdapter) LoadSettings(i *types.Instance) error {
	settingsPath := a.GetFilePath(i.UUID, InstanceSettingsPath)

	settingsBytes, err := os.ReadFile(settingsPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Warn("settings file not found. using default.")
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

func (a *InstanceFSAdapter) SaveService(instance *types.Instance) error {
	dir := a.GetFilePath(instance.UUID, InstanceVertexPath)
	servicePath := a.GetFilePath(instance.UUID, InstanceServicePath)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	serviceBytes, err := yaml.Marshal(instance.Service)
	if err != nil {
		return err
	}

	err = os.WriteFile(servicePath, serviceBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (a *InstanceFSAdapter) LoadService(instancePath string) (types.Service, error) {
	data, err := os.ReadFile(path.Join(instancePath, string(InstanceServicePath)))
	if err != nil {
		log.Warn("file not found",
			vlog.String("path", path.Dir(instancePath)),
		)
	}

	var service types.Service
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *InstanceFSAdapter) LoadServiceRaw(instancePath string) (interface{}, error) {
	servicePath := path.Join(instancePath, string(InstanceServicePath))

	data, err := os.ReadFile(servicePath)
	if err != nil {
		return nil, err
	}

	var service interface{}
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *InstanceFSAdapter) SaveEnv(i *types.Instance) error {
	envPath := a.GetFilePath(i.UUID, InstanceEnvPath)

	file, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	for key, value := range i.Env {
		_, err := file.WriteString(strings.Join([]string{key, value}, "=") + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *InstanceFSAdapter) LoadEnv(i *types.Instance) error {
	envPath := a.GetFilePath(i.UUID, InstanceEnvPath)

	file, err := os.Open(envPath)
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

		i.Env[line[0]] = line[1]
	}

	return nil
}

func (a *InstanceFSAdapter) Reload(load func(uuid uuid.UUID)) {
	entries, err := os.ReadDir(a.instancesPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Error(err)
			continue
		}

		isInstance := entry.IsDir() || info.Mode()&os.ModeSymlink != 0

		if isInstance {
			log.Info("found instance",
				vlog.String("uuid", entry.Name()),
			)

			id, err := uuid.Parse(entry.Name())
			if err != nil {
				log.Error(err)
				continue
			}

			load(id)
		}
	}
}

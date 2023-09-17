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
	"gopkg.in/yaml.v2"
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
	settingsPath := path.Join(a.GetPath(i.UUID), ".vertex", "instance_settings.json")

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
	settingsPath := path.Join(a.GetPath(i.UUID), ".vertex", "instance_settings.json")
	settingsBytes, err := os.ReadFile(settingsPath)

	if errors.Is(err, os.ErrNotExist) {
		log.Warn("instance_settings.json not found. using default.")
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

func (a *InstanceFSAdapter) ReadService(instancePath string) (types.Service, error) {
	data, err := os.ReadFile(path.Join(instancePath, ".vertex", "service.yml"))
	if err != nil {
		log.Warn("'.vertex/service.yml' file not found",
			vlog.String("path", path.Dir(instancePath)),
		)
	}

	var service types.Service
	err = yaml.Unmarshal(data, &service)
	return service, err
}

func (a *InstanceFSAdapter) SaveEnv(i *types.Instance, variables map[string]string) error {
	filepath := path.Join(a.GetPath(i.UUID), ".env")

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

func (a *InstanceFSAdapter) LoadEnv(i *types.Instance) error {
	filepath := path.Join(a.GetPath(i.UUID), ".env")

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

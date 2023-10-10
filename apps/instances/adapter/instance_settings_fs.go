package adapter

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
)

const InstanceSettingsPath = ".vertex/instance_settings.json"

type InstanceSettingsFSAdapter struct {
	instancesPath string
}

type InstanceSettingsFSAdapterParams struct {
	instancesPath string
}

func NewInstanceSettingsFSAdapter(params *InstanceSettingsFSAdapterParams) types.InstanceSettingsAdapterPort {
	if params == nil {
		params = &InstanceSettingsFSAdapterParams{}
	}
	if params.instancesPath == "" {
		params.instancesPath = path.Join(storage.Path, "instances")
	}

	adapter := &InstanceSettingsFSAdapter{
		instancesPath: params.instancesPath,
	}

	return adapter
}

func (a *InstanceSettingsFSAdapter) Save(uuid uuid.UUID, settings types.InstanceSettings) error {
	settingsPath := path.Join(a.instancesPath, uuid.String(), InstanceSettingsPath)

	settingsBytes, err := json.MarshalIndent(settings, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(settingsPath, settingsBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (a *InstanceSettingsFSAdapter) Load(uuid uuid.UUID) (types.InstanceSettings, error) {
	settingsPath := path.Join(a.instancesPath, uuid.String(), InstanceSettingsPath)

	settingsBytes, err := os.ReadFile(settingsPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Warn("settings file not found. using default.")
		return types.InstanceSettings{}, nil
	} else if err != nil {
		return types.InstanceSettings{}, err
	}

	var settings types.InstanceSettings
	err = json.Unmarshal(settingsBytes, &settings)
	return settings, err
}

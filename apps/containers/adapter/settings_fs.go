package adapter

import (
	"errors"
	"os"
	"path"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"gopkg.in/yaml.v3"
)

const ContainerSettingsPath = ".vertex/settings.yml"

type settingsFSAdapter struct {
	containersPath string
}

type SettingsFSAdapterParams struct {
	containersPath string
}

func NewSettingsFSAdapter(params *SettingsFSAdapterParams) port.SettingsAdapter {
	if params == nil {
		params = &SettingsFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	adapter := &settingsFSAdapter{
		containersPath: params.containersPath,
	}

	return adapter
}

func (a *settingsFSAdapter) Save(uuid types.ContainerID, settings types.ContainerSettings) error {
	settingsPath := path.Join(a.containersPath, uuid.String(), ContainerSettingsPath)

	settingsBytes, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	err = os.WriteFile(settingsPath, settingsBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (a *settingsFSAdapter) Load(uuid types.ContainerID) (types.ContainerSettings, error) {
	settingsPath := path.Join(a.containersPath, uuid.String(), ContainerSettingsPath)

	settingsBytes, err := os.ReadFile(settingsPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Warn("settings file not found. using default.")
		return types.ContainerSettings{}, nil
	} else if err != nil {
		return types.ContainerSettings{}, err
	}

	var settings types.ContainerSettings
	err = yaml.Unmarshal(settingsBytes, &settings)
	return settings, err
}

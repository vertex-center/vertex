package adapter

import (
	"errors"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"gopkg.in/yaml.v3"
)

const ContainerSettingsPath = ".vertex/settings.yml"

type ContainerSettingsFSAdapter struct {
	containersPath string
}

type ContainerSettingsFSAdapterParams struct {
	containersPath string
}

func NewContainerSettingsFSAdapter(params *ContainerSettingsFSAdapterParams) port.ContainerSettingsAdapter {
	if params == nil {
		params = &ContainerSettingsFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	adapter := &ContainerSettingsFSAdapter{
		containersPath: params.containersPath,
	}

	return adapter
}

func (a *ContainerSettingsFSAdapter) Save(uuid uuid.UUID, settings types.ContainerSettings) error {
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

func (a *ContainerSettingsFSAdapter) Load(uuid uuid.UUID) (types.ContainerSettings, error) {
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

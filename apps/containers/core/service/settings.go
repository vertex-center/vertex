package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type settingsService struct {
	adapter port.SettingsAdapter
}

func NewSettingsService(adapter port.SettingsAdapter) port.SettingsService {
	return &settingsService{
		adapter: adapter,
	}
}

func (s *settingsService) Save(inst *types.Container, settings types.ContainerSettings) error {
	inst.ContainerSettings = settings
	return s.adapter.Save(inst.UUID, settings)
}

func (s *settingsService) Load(inst *types.Container) error {
	settings, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	if settings.DisplayName == "" {
		settings.DisplayName = inst.Service.Name
	}
	inst.ContainerSettings = settings
	return nil
}

func (s *settingsService) SetLaunchOnStartup(inst *types.Container, value bool) error {
	inst.ContainerSettings.LaunchOnStartup = &value
	return s.adapter.Save(inst.UUID, inst.ContainerSettings)
}

func (s *settingsService) SetDisplayName(inst *types.Container, value string) error {
	inst.ContainerSettings.DisplayName = value
	return s.adapter.Save(inst.UUID, inst.ContainerSettings)
}

func (s *settingsService) SetDatabases(inst *types.Container, databases map[string]types.ContainerID) error {
	inst.Databases = databases
	return s.adapter.Save(inst.UUID, inst.ContainerSettings)
}

func (s *settingsService) SetVersion(inst *types.Container, value string) error {
	inst.Version = &value
	return s.adapter.Save(inst.UUID, inst.ContainerSettings)
}

func (s *settingsService) SetTags(inst *types.Container, tags []string) error {
	inst.Tags = tags
	return s.adapter.Save(inst.UUID, inst.ContainerSettings)
}

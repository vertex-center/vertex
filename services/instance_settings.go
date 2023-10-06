package services

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type InstanceSettingsService struct {
	adapter types.InstanceSettingsAdapterPort
}

func NewInstanceSettingsService(instanceSettingsAdapter types.InstanceSettingsAdapterPort) InstanceSettingsService {
	return InstanceSettingsService{
		adapter: instanceSettingsAdapter,
	}
}

func (s *InstanceSettingsService) Save(inst *types.Instance, settings types.InstanceSettings) error {
	inst.InstanceSettings = settings
	return s.adapter.Save(inst.UUID, settings)
}

func (s *InstanceSettingsService) Load(inst *types.Instance) error {
	settings, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	inst.InstanceSettings = settings
	return nil
}

func (s *InstanceSettingsService) SetLaunchOnStartup(inst *types.Instance, value bool) error {
	inst.InstanceSettings.LaunchOnStartup = &value
	return s.adapter.Save(inst.UUID, inst.InstanceSettings)
}

func (s *InstanceSettingsService) SetDisplayName(inst *types.Instance, value string) error {
	inst.InstanceSettings.DisplayName = &value
	return s.adapter.Save(inst.UUID, inst.InstanceSettings)
}

func (s *InstanceSettingsService) SetDatabases(inst *types.Instance, databases map[string]uuid.UUID) error {
	inst.Databases = databases
	return s.adapter.Save(inst.UUID, inst.InstanceSettings)
}

func (s *InstanceSettingsService) SetVersion(inst *types.Instance, value string) error {
	inst.Version = &value
	return s.adapter.Save(inst.UUID, inst.InstanceSettings)
}

func (s *InstanceSettingsService) SetTags(inst *types.Instance, tags []string) error {
	inst.Tags = tags
	return s.adapter.Save(inst.UUID, inst.InstanceSettings)
}

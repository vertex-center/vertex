package service

import (
	"errors"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
)

type SettingsService struct {
	adapter port.SettingsAdapter
}

func NewSettingsService(adapter port.SettingsAdapter) port.SettingsService {
	return &SettingsService{
		adapter: adapter,
	}
}

func (s *SettingsService) Get() (types.AdminSettings, error) {
	return s.adapter.Get()
}

func (s *SettingsService) Update(settings types.AdminSettings) error {
	var errs []error
	if settings.Webhook != nil {
		if err := s.SetWebhook(*settings.Webhook); err != nil {
			errs = append(errs, err)
		}
	}
	if settings.UpdatesChannel != "" {
		if err := s.SetChannel(settings.UpdatesChannel); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (s *SettingsService) GetWebhook() (*string, error) {
	settings, err := s.Get()
	if err != nil {
		return nil, err
	}
	return settings.Webhook, nil
}

func (s *SettingsService) SetWebhook(webhook string) error {
	return s.adapter.SetWebhook(webhook)
}

func (s *SettingsService) GetChannel() (types.UpdatesChannel, error) {
	settings, err := s.Get()
	if err != nil {
		return types.UpdatesChannelStable, err
	}
	return settings.UpdatesChannel, nil
}

func (s *SettingsService) SetChannel(channel types.UpdatesChannel) error {
	return s.adapter.SetChannel(channel)
}

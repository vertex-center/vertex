package service

import (
	"errors"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
)

type settingsService struct {
	adapter port.SettingsAdapter
}

func NewSettingsService(adapter port.SettingsAdapter) port.SettingsService {
	return &settingsService{
		adapter: adapter,
	}
}

func (s *settingsService) Get() (types.AdminSettings, error) {
	return s.adapter.Get()
}

func (s *settingsService) Update(settings types.AdminSettings) error {
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

func (s *settingsService) GetWebhook() (*string, error) {
	settings, err := s.Get()
	if err != nil {
		return nil, err
	}
	return settings.Webhook, nil
}

func (s *settingsService) SetWebhook(webhook string) error {
	return s.adapter.SetWebhook(webhook)
}

func (s *settingsService) GetChannel() (types.UpdatesChannel, error) {
	settings, err := s.Get()
	if err != nil {
		return types.UpdatesChannelStable, err
	}
	return settings.UpdatesChannel, nil
}

func (s *settingsService) SetChannel(channel types.UpdatesChannel) error {
	return s.adapter.SetChannel(channel)
}

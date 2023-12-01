package service

import (
	"errors"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
)

type AdminSettingsService struct {
	adapter port.AdminSettingsAdapter
}

func NewAdminSettingsService(adapter port.AdminSettingsAdapter) port.AdminSettingsService {
	return &AdminSettingsService{
		adapter: adapter,
	}
}

func (s *AdminSettingsService) Get() (types.AdminSettings, error) {
	return s.adapter.Get()
}

func (s *AdminSettingsService) Update(settings types.AdminSettings) error {
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

func (s *AdminSettingsService) GetWebhook() (*string, error) {
	settings, err := s.Get()
	if err != nil {
		return nil, err
	}
	return settings.Webhook, nil
}

func (s *AdminSettingsService) SetWebhook(webhook string) error {
	return s.adapter.SetWebhook(webhook)
}

func (s *AdminSettingsService) GetChannel() (types.UpdatesChannel, error) {
	settings, err := s.Get()
	if err != nil {
		return types.UpdatesChannelStable, err
	}
	return settings.UpdatesChannel, nil
}

func (s *AdminSettingsService) SetChannel(channel types.UpdatesChannel) error {
	return s.adapter.SetChannel(channel)
}

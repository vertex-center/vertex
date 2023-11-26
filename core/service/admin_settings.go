package service

import (
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
	return s.adapter.Update(settings)
}

func (s *AdminSettingsService) GetWebhook() (*string, error) {
	settings, err := s.Get()
	if err != nil {
		return nil, err
	}
	return settings.Webhook, nil
}

func (s *AdminSettingsService) SetWebhook(webhook string) error {
	return s.Update(types.AdminSettings{
		Webhook: &webhook,
	})
}

func (s *AdminSettingsService) GetChannel() (types.UpdatesChannel, error) {
	settings, err := s.Get()
	if err != nil {
		return types.UpdatesChannelStable, err
	}
	if settings.UpdatesChannel == nil {
		return types.UpdatesChannelStable, nil
	}
	return *settings.UpdatesChannel, nil
}

func (s *AdminSettingsService) SetChannel(channel types.UpdatesChannel) error {
	return s.Update(types.AdminSettings{
		UpdatesChannel: &channel,
	})
}

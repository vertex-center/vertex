package services

import (
	"github.com/vertex-center/vertex/types"
)

type SettingsService struct {
	settingsAdapter types.SettingsAdapterPort
}

func NewSettingsService(settingsAdapter types.SettingsAdapterPort) SettingsService {
	return SettingsService{
		settingsAdapter: settingsAdapter,
	}
}

func (s *SettingsService) Update(settings types.Settings) error {
	if settings.Notifications != nil {
		notifs := settings.Notifications
		if notifs.Webhook != nil {
			err := s.settingsAdapter.SetNotificationsWebhook(notifs.Webhook)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SettingsService) GetNotificationsWebhook() *string {
	return s.settingsAdapter.GetNotificationsWebhook()
}

func (s *SettingsService) SetNotificationsWebhook(webhook *string) error {
	return s.settingsAdapter.SetNotificationsWebhook(webhook)
}

package services

import (
	"github.com/vertex-center/vertex/types"
)

type SettingsService struct {
	settingsRepo types.SettingsRepository
}

func NewSettingsService(settingsRepo types.SettingsRepository) SettingsService {
	return SettingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *SettingsService) Update(settings types.Settings) error {
	if settings.Notifications != nil {
		notifs := settings.Notifications
		if notifs.Webhook != nil {
			err := s.settingsRepo.SetNotificationsWebhook(notifs.Webhook)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SettingsService) GetNotificationsWebhook() *string {
	return s.settingsRepo.GetNotificationsWebhook()
}

func (s *SettingsService) SetNotificationsWebhook(webhook *string) error {
	return s.settingsRepo.SetNotificationsWebhook(webhook)
}

package types

type SettingsNotifications struct {
	Webhook *string `json:"webhook,omitempty"`
}

type Settings struct {
	Notifications *SettingsNotifications `json:"notifications,omitempty"`
}

type SettingsRepository interface {
	GetSettings() Settings

	GetNotificationsWebhook() *string
	SetNotificationsWebhook(webhook *string) error
}

package types

type SettingsNotifications struct {
	Webhook *string `json:"webhook,omitempty"`
}

type SettingsUpdatesChannel string

const (
	SettingsUpdatesChannelStable SettingsUpdatesChannel = "stable"
	SettingsUpdatesChannelBeta   SettingsUpdatesChannel = "beta"
)

type SettingsUpdates struct {
	Channel *SettingsUpdatesChannel `json:"channel,omitempty"`
}

type Settings struct {
	Notifications *SettingsNotifications `json:"notifications,omitempty"`
	Updates       *SettingsUpdates       `json:"updates,omitempty"`
}

type SettingsAdapterPort interface {
	GetSettings() Settings

	GetNotificationsWebhook() *string
	SetNotificationsWebhook(webhook *string) error

	GetChannel() *SettingsUpdatesChannel
	SetChannel(channel *SettingsUpdatesChannel) error
}

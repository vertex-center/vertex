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

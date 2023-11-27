package types

type UpdatesChannel string

const (
	UpdatesChannelStable UpdatesChannel = "stable"
	UpdatesChannelBeta   UpdatesChannel = "beta"
)

type AdminSettings struct {
	UpdatesChannel UpdatesChannel `json:"updates_channel,omitempty" gorm:"default:'stable'"`
	Webhook        *string        `json:"webhook,omitempty"`
}

func NewAdminSettings() AdminSettings {
	return AdminSettings{
		UpdatesChannel: UpdatesChannelStable,
	}
}

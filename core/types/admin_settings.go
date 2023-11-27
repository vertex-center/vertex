package types

type UpdatesChannel string

const (
	UpdatesChannelStable UpdatesChannel = "stable"
	UpdatesChannelBeta   UpdatesChannel = "beta"
)

type AdminSettings struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement:false;default:1;check:id=1"`
	UpdatesChannel UpdatesChannel `json:"updates_channel,omitempty" gorm:"default:'stable';check:updates_channel IN ('stable', 'beta')"`
	Webhook        *string        `json:"webhook,omitempty"`
}

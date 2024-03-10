package types

import "github.com/vertex-center/vertex/server/common/baseline"

type AdminSettings struct {
	ID             uint             `json:"id" db:"id"`
	UpdatesChannel baseline.Channel `json:"updates_channel,omitempty" db:"updates_channel"`
	Webhook        *string          `json:"webhook,omitempty" db:"webhook"`
	CreatedAt      int64            `json:"created_at" db:"created_at"`
	UpdatedAt      int64            `json:"updated_at" db:"updated_at"`
	DeletedAt      *int64           `json:"deleted_at,omitempty" db:"deleted_at"`
}

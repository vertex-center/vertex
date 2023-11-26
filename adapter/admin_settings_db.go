package adapter

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
)

type AdminSettingsDbAdapter struct {
	db port.DbConfigAdapter
}

func NewAdminSettingsDbAdapter(db port.DbConfigAdapter) port.AdminSettingsAdapter {
	return &AdminSettingsDbAdapter{
		db: db,
	}
}

func (s AdminSettingsDbAdapter) Get() (types.AdminSettings, error) {
	var settings types.AdminSettings
	err := s.db.Get().First(&settings).Error
	return settings, err
}

func (s AdminSettingsDbAdapter) Update(settings types.AdminSettings) error {
	return s.db.Get().Save(&settings).Error
}

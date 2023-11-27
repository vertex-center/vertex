package adapter

import (
	"errors"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"gorm.io/gorm"
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return types.NewAdminSettings(), nil
	}
	return settings, err
}

func (s AdminSettingsDbAdapter) Update(settings types.AdminSettings) error {
	return s.db.Get().Save(&settings).Error
}

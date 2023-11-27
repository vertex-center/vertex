package adapter

import (
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
	settings.ID = 1
	err := s.db.Get().Debug().FirstOrCreate(&settings).Error
	return settings, err
}

func (s AdminSettingsDbAdapter) Update(settings types.AdminSettings) error {
	return s.db.Get().Debug().Transaction(func(tx *gorm.DB) error {
		var current types.AdminSettings
		current.ID = 1
		err := tx.FirstOrCreate(&current).Error
		if err != nil {
			return err
		}
		settings.ID = 1
		return tx.Model(&current).Updates(&settings).Error
	})
}

package adapter

import (
	"sync"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"gorm.io/gorm"
)

type AdminSettingsDbAdapter struct {
	db port.DbAdapter
	// dbMutex is used to prevent concurrent access to the settings table.
	// This is also needed to ensure that there is only one row in the table.
	dbMutex sync.RWMutex
}

func NewAdminSettingsDbAdapter(db port.DbAdapter) port.AdminSettingsAdapter {
	return &AdminSettingsDbAdapter{
		db: db,
	}
}

func (s *AdminSettingsDbAdapter) Get() (types.AdminSettings, error) {
	s.dbMutex.Lock()
	defer s.dbMutex.Unlock()

	var settings types.AdminSettings
	err := s.db.Get().FirstOrCreate(&settings).Error
	return settings, err
}

func (s *AdminSettingsDbAdapter) Update(settings types.AdminSettings) error {
	s.dbMutex.Lock()
	defer s.dbMutex.Unlock()

	return s.db.Get().Transaction(func(tx *gorm.DB) error {
		var current types.AdminSettings
		err := tx.FirstOrCreate(&current).Error
		if err != nil {
			return err
		}
		return tx.Model(&current).Updates(&settings).Error
	})
}

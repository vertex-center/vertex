package adapter

import (
	"sync"
	"time"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/storage"
)

type settingsDbAdapter struct {
	db storage.DB
	// dbMutex is used to prevent concurrent access to the settings table.
	// This is also needed to ensure that there is only one row in the table.
	dbMutex sync.RWMutex
}

func NewSettingsDbAdapter(db storage.DB) port.SettingsAdapter {
	return &settingsDbAdapter{
		db: db,
	}
}

func (a *settingsDbAdapter) Get() (types.AdminSettings, error) {
	a.dbMutex.Lock()
	defer a.dbMutex.Unlock()

	var settings types.AdminSettings
	err := a.db.Get(&settings, "SELECT * FROM admin_settings LIMIT 1")
	if err != nil {
		return types.AdminSettings{}, err
	}
	return settings, err
}

func (a *settingsDbAdapter) SetChannel(channel types.UpdatesChannel) error {
	a.dbMutex.Lock()
	defer a.dbMutex.Unlock()

	_, err := a.db.Exec(`
		UPDATE admin_settings
		SET updates_channel = $1, updated_at = $2
		WHERE id = 1
	`, channel, time.Now().Unix())
	return err
}

func (a *settingsDbAdapter) SetWebhook(webhook string) error {
	a.dbMutex.Lock()
	defer a.dbMutex.Unlock()

	_, err := a.db.Exec(`
		UPDATE admin_settings
		SET webhook = $1, updated_at = $2
		WHERE id = 1
	`, webhook, time.Now().Unix())
	return err
}

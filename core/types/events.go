package types

import "gorm.io/gorm"

type (
	// EventServerStart is dispatched when the server is started. This event is
	// dispatched before the setup.
	EventServerStart struct{}

	// EventServerSetupCompleted is dispatched when the server setup is completed.
	EventServerSetupCompleted struct{}

	// EventAppReady is dispatched when the app with the id AppID is ready to be used.
	EventAppReady struct {
		// AppID is the id of the app that is ready.
		AppID string
	}

	// EventAllAppsReady is dispatched when all apps are ready to be used.
	EventAllAppsReady struct{}

	// EventDbMigrate is dispatched when the database is migrated.
	// Use this event to migrate the database of your app.
	EventDbMigrate struct {
		Db *gorm.DB
	}

	// EventDbCopy is dispatched when the database is copied.
	// Use this event to send which tables you want to copy to the new database.
	EventDbCopy struct {
		tables *[]interface{}
	}

	// EventServerStop is dispatched when the server is stopped.
	EventServerStop struct{}

	// EventServerHardReset is dispatched when the server is hard reset. This is used for testing purposes.
	EventServerHardReset struct{}

	// EventVertexUpdated is dispatched when the vertex binary is updated.
	EventVertexUpdated struct{}
)

// AddTable adds a table to the list of tables that will be copied to the new database.
// Example usage: e.AddTable(types.User{})
func (e *EventDbCopy) AddTable(t ...interface{}) {
	if e.tables == nil {
		e.tables = &[]interface{}{}
	}
	*e.tables = append(*e.tables, t...)
}

// All returns all tables that will be copied to the new database.
func (e *EventDbCopy) All() []interface{} {
	if e.tables == nil {
		return []interface{}{}
	}
	return *e.tables
}

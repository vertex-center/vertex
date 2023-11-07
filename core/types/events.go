package types

type (
	// EventServerStart is dispatched when the server is started.
	EventServerStart struct {
		// PostMigrationCommands are commands that should be executed after the server has started.
		// These are migration commands that cannot be executed before the server has started.
		PostMigrationCommands []interface{}
	}

	// EventAppReady is dispatched when the app with the id AppID is ready to be used.
	EventAppReady struct {
		// AppID is the id of the app that is ready.
		AppID string
	}

	// EventServerStop is dispatched when the server is stopped.
	EventServerStop struct{}

	// EventServerHardReset is dispatched when the server is hard reset. This is used for testing purposes.
	EventServerHardReset struct{}

	// EventVertexUpdated is dispatched when the vertex binary is updated.
	EventVertexUpdated struct{}
)

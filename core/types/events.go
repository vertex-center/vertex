package types

type (
	EventServerStart struct {
		// PostMigrationCommands are commands that should be executed after the server has started.
		// These are migration commands that cannot be executed before the server has started.
		PostMigrationCommands []interface{}
	}

	EventAppReady struct {
		AppID string
	}

	EventServerStop      struct{}
	EventServerHardReset struct{}
	EventVertexUpdated   struct{}
)

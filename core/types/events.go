package types

type (
	// EventServerDownloadDependencies is dispatched before the server starts.
	EventServerDownloadDependencies struct{}

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

	// EventServerStop is dispatched when the server is stopped.
	EventServerStop struct{}

	// EventServerHardReset is dispatched when the server is hard reset. This is used for testing purposes.
	EventServerHardReset struct{}

	// EventVertexUpdated is dispatched when the vertex binary is updated.
	EventVertexUpdated struct{}
)

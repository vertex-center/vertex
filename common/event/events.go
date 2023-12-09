package event

type (
	// ServerLoad is dispatched before the server starts.
	ServerLoad struct{}

	// ServerStart is dispatched when the server is started. This event is
	// dispatched before the setup.
	ServerStart struct{}

	// ServerSetupCompleted is dispatched when the server setup is completed.
	ServerSetupCompleted struct{}

	// AppReady is dispatched when the app with the id AppID is ready to be used.
	AppReady struct {
		// AppID is the id of the app that is ready.
		AppID string
	}

	// AllAppsReady is dispatched when all apps are ready to be used.
	AllAppsReady struct{}

	// ServerStop is dispatched when the server is stopped.
	ServerStop struct{}

	// ServerHardReset is dispatched when the server is hard reset. This is used for testing purposes.
	ServerHardReset struct{}

	// VertexUpdated is dispatched when the vertex binary is updated.
	VertexUpdated struct{}
)

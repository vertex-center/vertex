package event

type (
	ServerLoad           struct{} // ServerLoad is dispatched before the server starts.
	ServerStart          struct{} // ServerStart is dispatched when the server is started.
	ServerSetupCompleted struct{} // ServerSetupCompleted is dispatched when the server setup is completed.
	ServerStop           struct{} // ServerStop is dispatched when the server is stopped.
	VertexUpdated        struct{} // VertexUpdated is dispatched when the vertex binary is updated.
)

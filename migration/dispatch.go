package migration

// CommandRecreateContainers allows a migration to recreate all containers.
type CommandRecreateContainers struct{}

// The CommandsDispatcher allows a migration to dispatch commands to Vertex.
// e.g. Send CommandRecreateContainers if a migration needs to recreate all containers.
type CommandsDispatcher interface {
	DispatchCommands() []interface{}
}

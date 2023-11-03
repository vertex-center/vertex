package event

// Bus is the interface that must be implemented by an event bus.
type Bus interface {
	// AddListener allows a listener to listen to events
	// dispatched on this Bus.
	AddListener(l Listener)
	// RemoveListener allows a listener to stop listening
	// to events dispatched on this Bus.
	RemoveListener(l Listener)
	// DispatchEvent dispatches an event to all listeners
	// listening to this Bus.
	DispatchEvent(e Event)
}

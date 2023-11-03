package event

// EventBus is the interface that must be implemented by an event bus.
type EventBus interface {
	// AddListener allows a listener to listen to events
	// dispatched on this EventBus.
	AddListener(l EventListener)
	// RemoveListener allows a listener to stop listening
	// to events dispatched on this EventBus.
	RemoveListener(l EventListener)
	// DispatchEvent dispatches an event to all listeners
	// listening to this EventBus.
	DispatchEvent(e Event)
}

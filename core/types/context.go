package types

type VertexContext struct {
	eventBus *EventBus
}

func NewVertexContext() *VertexContext {
	return &VertexContext{
		eventBus: NewEventBus(),
	}
}

func (c *VertexContext) DispatchEvent(e interface{}) {
	c.eventBus.Send(e)
}

func (c *VertexContext) AddListener(l Listener) {
	c.eventBus.AddListener(l)
}

func (c *VertexContext) RemoveListener(l Listener) {
	c.eventBus.RemoveListener(l)
}

package common

import (
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
)

type VertexContext struct {
	about  About
	bus    event.Bus
	kernel bool
}

func NewVertexContext(about About, kernel bool) *VertexContext {
	return &VertexContext{
		about:  about,
		bus:    event.NewMemoryBus(),
		kernel: kernel,
	}
}

func (c *VertexContext) DispatchEvent(e event.Event) {
	err := c.DispatchEventWithErr(e)
	if err != nil {
		log.Error(err)
	}
}

func (c *VertexContext) DispatchEventWithErr(e event.Event) error {
	return c.bus.DispatchEvent(e)
}

func (c *VertexContext) AddListener(l event.Listener) {
	c.bus.AddListener(l)
}

func (c *VertexContext) RemoveListener(l event.Listener) {
	c.bus.RemoveListener(l)
}

func (c *VertexContext) About() About {
	return c.about
}

func (c *VertexContext) Kernel() bool {
	return c.kernel
}

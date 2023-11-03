package types

import (
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/event/types"
	"github.com/vertex-center/vertex/pkg/log"
)

type VertexContext struct {
	bus types.EventBus
}

func NewVertexContext() *VertexContext {
	return &VertexContext{
		bus: event.NewMemoryBus(),
	}
}

func (c *VertexContext) DispatchEvent(e interface{}) {
	if _, ok := e.(EventServerHardReset); ok {
		if !config.Current.Debug() {
			log.Warn("hard reset event received but skipped; this can be a malicious application, or you may have forgotten to switch to the development mode.")
			return
		}
		log.Warn("hard reset event dispatched.")
	}

	c.bus.DispatchEvent(e)
}

func (c *VertexContext) AddListener(l types.EventListener) {
	c.bus.AddListener(l)
}

func (c *VertexContext) RemoveListener(l types.EventListener) {
	c.bus.RemoveListener(l)
}

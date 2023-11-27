package types

import (
	"reflect"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type VertexContext struct {
	bus event.Bus
}

func NewVertexContext() *VertexContext {
	return &VertexContext{
		bus: event.NewMemoryBus(),
	}
}

func (c *VertexContext) DispatchEvent(e event.Event) {
	log.Debug("dispatching event", vlog.String("name", reflect.TypeOf(e).String()))

	if _, ok := e.(EventServerHardReset); ok {
		if !config.Current.Debug() {
			log.Warn("hard reset event received but skipped; this can be a malicious application, or you may have forgotten to switch to the development mode.")
			return
		}
		log.Warn("hard reset event dispatched.")
	}

	c.bus.DispatchEvent(e)
}

func (c *VertexContext) AddListener(l event.Listener) {
	c.bus.AddListener(l)
}

func (c *VertexContext) RemoveListener(l event.Listener) {
	c.bus.RemoveListener(l)
}

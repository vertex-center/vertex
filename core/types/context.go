package types

import (
	"reflect"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type VertexContext struct {
	about  About
	bus    event.Bus
	db     *DB
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
	log.Debug("dispatching event", vlog.String("name", reflect.TypeOf(e).String()))

	if _, ok := e.(EventServerHardReset); ok {
		if !config.Current.Debug() {
			log.Warn("hard reset event received but skipped; this can be a malicious application, or you may have forgotten to switch to the development mode.")
			return nil
		}
		log.Warn("hard reset event dispatched.")
	}

	return c.bus.DispatchEvent(e)
}

func (c *VertexContext) AddListener(l event.Listener) {
	c.bus.AddListener(l)
}

func (c *VertexContext) RemoveListener(l event.Listener) {
	c.bus.RemoveListener(l)
}

func (c *VertexContext) SetDb(db *DB) {
	c.db = db
}

func (c *VertexContext) Db() *DB {
	return c.db
}

func (c *VertexContext) About() About {
	return c.about
}

func (c *VertexContext) Kernel() bool {
	return c.kernel
}

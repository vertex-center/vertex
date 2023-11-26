package service

import (
	"context"
	"errors"

	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
)

func (s *SetupService) startContainer(inst *types.Container) error {
	eventsChan := make(chan interface{})
	defer close(eventsChan)

	abortChan := make(chan bool)
	defer close(abortChan)

	l := event.NewTempListener(func(e event.Event) {
		switch e := e.(type) {
		case types.EventContainerStatusChange:
			if e.ContainerUUID != inst.UUID {
				return
			}
			eventsChan <- e
		}
	})

	s.ctx.AddListener(l)
	defer s.ctx.RemoveListener(l)

	client := containersapi.NewContainersClient()

	go func() {
		apiError := client.StartContainer(context.Background(), inst.UUID)
		if apiError != nil {
			log.Error(apiError.RouterError())
		}
		abortChan <- true
	}()

	errFailedToStart := errors.New("failed to start container")

	for {
		select {
		case e := <-eventsChan:
			switch e := e.(type) {
			case types.EventContainerStatusChange:
				if e.Status == types.ContainerStatusRunning {
					return nil
				} else if e.Status == types.ContainerStatusOff || e.Status == types.ContainerStatusError {
					return errFailedToStart
				}
			}
		case <-abortChan:
			return errFailedToStart
		}
	}
}

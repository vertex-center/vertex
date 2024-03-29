package service

import (
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/pkg/event"
)

func (s *sqlService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *sqlService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case types.EventContainerStatusChange:
		if e.Status == types.ContainerStatusRunning {
			s.onContainerStart(&e.Container)
		} else if e.Status == types.ContainerStatusOff {
			s.onContainerStop(e.ContainerID)
		}
	}
	return nil
}

func (s *sqlService) onContainerStart(inst *types.Container) {
	_, err := s.getDbFeature(inst)
	if err != nil {
		// Not a SQL database
		return
	}

	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[inst.ID]; ok {
		return
	}

	dbms, err := s.createDbmsAdapter(inst)
	if err != nil {
		return
	}

	s.dbms[inst.ID] = dbms
}

func (s *sqlService) onContainerStop(uuid uuid.UUID) {
	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[uuid]; !ok {
		return
	}

	delete(s.dbms, uuid)
}

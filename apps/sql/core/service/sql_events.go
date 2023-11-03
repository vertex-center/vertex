package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	evtypes "github.com/vertex-center/vertex/pkg/event/types"
)

func (s *SqlService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *SqlService) OnEvent(e evtypes.Event) {
	switch e := e.(type) {
	case types.EventContainerStatusChange:
		if e.Status == types.ContainerStatusRunning {
			s.onContainerStart(&e.Container)
		} else if e.Status == types.ContainerStatusOff {
			s.onContainerStop(e.ContainerUUID)
		}
	}
}

func (s *SqlService) onContainerStart(inst *types.Container) {
	_, err := s.getDbFeature(inst)
	if err != nil {
		// Not a SQL database
		return
	}

	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[inst.UUID]; ok {
		return
	}

	dbms, err := s.createDbmsAdapter(inst)
	if err != nil {
		return
	}

	s.dbms[inst.UUID] = dbms
}

func (s *SqlService) onContainerStop(uuid uuid.UUID) {
	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[uuid]; !ok {
		return
	}

	delete(s.dbms, uuid)
}

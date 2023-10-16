package service

import (
	"github.com/google/uuid"
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"
)

func (s *SqlService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *SqlService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types2.EventContainerStatusChange:
		if e.Status == types2.ContainerStatusRunning {
			s.onContainerStart(&e.Container)
		} else if e.Status == types2.ContainerStatusOff {
			s.onContainerStop(e.ContainerUUID)
		}
	}
}

func (s *SqlService) onContainerStart(inst *types2.Container) {
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

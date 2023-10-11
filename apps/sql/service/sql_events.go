package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/types"
)

func (s *SqlService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *SqlService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceStatusChange:
		if e.Status == types.InstanceStatusRunning {
			s.onInstanceStart(&e.Instance)
		} else if e.Status == types.InstanceStatusOff {
			s.onInstanceStop(e.InstanceUUID)
		}
	}
}

func (s *SqlService) onInstanceStart(inst *types.Instance) {
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

func (s *SqlService) onInstanceStop(uuid uuid.UUID) {
	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[uuid]; !ok {
		return
	}

	delete(s.dbms, uuid)
}

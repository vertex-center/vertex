package service

import (
	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
)

func (s *SqlService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case instancestypes.EventInstanceStatusChange:
		if e.Status == instancestypes.InstanceStatusRunning {
			s.onInstanceStart(&e.Instance)
		} else if e.Status == instancestypes.InstanceStatusOff {
			s.onInstanceStop(e.InstanceUUID)
		}
	}
}

func (s *SqlService) onInstanceStart(inst *instancestypes.Instance) {
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

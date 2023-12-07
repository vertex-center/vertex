package service

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
)

type hardwareKernelService struct {
	adapter port.HardwareKernelAdapter
}

func NewHardwareKernelService(adapter port.HardwareKernelAdapter) port.HardwareKernelService {
	return &hardwareKernelService{
		adapter: adapter,
	}
}

func (s *hardwareKernelService) Reboot() error {
	return s.adapter.Reboot()
}

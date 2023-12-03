package service

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
)

type HardwareKernelService struct {
	adapter port.HardwareKernelAdapter
}

func NewHardwareKernelService(adapter port.HardwareKernelAdapter) port.HardwareKernelService {
	return &HardwareKernelService{
		adapter: adapter,
	}
}

func (s *HardwareKernelService) Reboot() error {
	return s.adapter.Reboot()
}

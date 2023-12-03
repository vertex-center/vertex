package adapter

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/reboot"
)

type HardwareKernelAdapter struct{}

func NewHardwareKernelAdapter() port.HardwareKernelAdapter {
	return HardwareKernelAdapter{}
}

func (HardwareKernelAdapter) Reboot() error {
	return reboot.Reboot()
}

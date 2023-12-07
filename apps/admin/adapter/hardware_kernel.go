package adapter

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/reboot"
)

type hardwareKernelAdapter struct{}

func NewHardwareKernelAdapter() port.HardwareKernelAdapter {
	return hardwareKernelAdapter{}
}

func (hardwareKernelAdapter) Reboot() error {
	return reboot.Reboot()
}

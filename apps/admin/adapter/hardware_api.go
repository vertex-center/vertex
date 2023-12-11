package adapter

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/api"
	"github.com/vertex-center/vertex/apps/admin/core/port"
)

type hardwareApiAdapter struct{}

func NewHardwareApiAdapter() port.HardwareAdapter {
	return hardwareApiAdapter{}
}

func (hardwareApiAdapter) Reboot(ctx context.Context) error {
	return api.NewAdminKernelClient(ctx).Reboot(ctx)
}

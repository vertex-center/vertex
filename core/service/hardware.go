package service

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/vdocker"
)

type HardwareService struct{}

func NewHardwareService() port.HardwareService {
	return &HardwareService{}
}

func (s HardwareService) Get() types.Hardware {
	stats, err := host.Info()
	if err != nil {
		// fallback to runtime.GOOS and runtime.GOARCH
		return types.Hardware{
			Dockerized: vdocker.RunningInDocker(),
			Host: types.Host{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
		}
	}

	return types.Hardware{
		Dockerized: vdocker.RunningInDocker(),
		Host: types.Host{
			OS:       stats.OS,
			Arch:     stats.KernelArch,
			Platform: stats.Platform,
			Version:  stats.PlatformVersion,
			Name:     stats.Hostname,
		},
	}
}

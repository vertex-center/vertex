package services

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/vertex-center/vertex/types"
)

type HardwareService struct{}

func NewHardwareService() HardwareService {
	return HardwareService{}
}

func (s HardwareService) GetHardware() types.Hardware {
	stats, err := host.Info()
	if err != nil {
		// fallback to runtime.GOOS and runtime.GOARCH
		return types.Hardware{
			Host: types.Host{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
		}
	}

	return types.Hardware{
		Host: types.Host{
			OS:       stats.OS,
			Arch:     stats.KernelArch,
			Platform: stats.Platform,
			Version:  stats.PlatformVersion,
			Name:     stats.Hostname,
		},
	}
}

package service

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
)

type HardwareService struct{}

func NewHardwareService() port.HardwareService {
	return &HardwareService{}
}

func (s HardwareService) GetHost() (types.Host, error) {
	info, err := host.Info()
	if err != nil {
		return types.Host{}, err
	}

	return types.Host{
		Hostname:             info.Hostname,
		Uptime:               info.Uptime,
		BootTime:             info.BootTime,
		Procs:                info.Procs,
		OS:                   info.OS,
		Platform:             info.Platform,
		PlatformFamily:       info.PlatformFamily,
		PlatformVersion:      info.PlatformVersion,
		KernelVersion:        info.KernelVersion,
		KernelArch:           info.KernelArch,
		VirtualizationSystem: info.VirtualizationSystem,
		VirtualizationRole:   info.VirtualizationRole,
		HostID:               info.HostID,
	}, nil
}

func (s HardwareService) GetCPUs() ([]types.CPU, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var cpus []types.CPU
	for _, i := range info {
		cpus = append(cpus, types.CPU{
			Count:      i.CPU,
			VendorID:   i.VendorID,
			Family:     i.Family,
			Model:      i.Model,
			Stepping:   i.Stepping,
			PhysicalID: i.PhysicalID,
			CoreID:     i.CoreID,
			CoresCount: i.Cores,
			ModelName:  i.ModelName,
			Mhz:        i.Mhz,
			CacheSize:  i.CacheSize,
			Flags:      i.Flags,
			Microcode:  i.Microcode,
		})
	}
	return cpus, nil
}

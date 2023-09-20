package services

import (
	"runtime"

	"github.com/vertex-center/vertex/types"
	"golang.org/x/sys/unix"
)

type HardwareService struct{}

func NewHardwareService() HardwareService {
	return HardwareService{}
}

func (s HardwareService) GetHardware() types.Hardware {
	uname := unix.Utsname{}
	err := unix.Uname(&uname)
	if err != nil {
		// fallback to runtime.GOOS and runtime.GOARCH
		return types.Hardware{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		}
	}

	return types.Hardware{
		OS:      unix.ByteSliceToString(uname.Sysname[:]),
		Arch:    unix.ByteSliceToString(uname.Machine[:]),
		Version: unix.ByteSliceToString(uname.Release[:]),
		Name:    unix.ByteSliceToString(uname.Nodename[:]),
	}
}

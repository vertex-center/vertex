package config

var KernelCurrent = NewKernel()

type Kernel struct {
	Config

	Uid uint32 `json:"uid"`
	Gid uint32 `json:"gid"`
}

func NewKernel() Kernel {
	return Kernel{
		Config: New(),
	}
}

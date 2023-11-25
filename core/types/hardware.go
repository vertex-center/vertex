// This package contains the types for the hardware service.
// The structs are almost identical to those of the `gopsutil` library,
// but they are copied here to avoid the dependency and to be able to
// add more fields.

package types

type Host struct {
	Hostname             string `json:"hostname,omitempty"`
	Uptime               uint64 `json:"uptime,omitempty"`
	BootTime             uint64 `json:"boot_time,omitempty"`
	Procs                uint64 `json:"procs,omitempty"`
	OS                   string `json:"os,omitempty"`
	Platform             string `json:"platform,omitempty"`
	PlatformFamily       string `json:"platform_family,omitempty"`
	PlatformVersion      string `json:"platform_version,omitempty"`
	KernelVersion        string `json:"kernel_version,omitempty"`
	KernelArch           string `json:"kernel_arch,omitempty"`
	VirtualizationSystem string `json:"virtualization_system,omitempty"`
	VirtualizationRole   string `json:"virtualization_role,omitempty"`
	HostID               string `json:"host_id,omitempty"`
}

type CPU struct {
	Count      int32    `json:"count,omitempty"`
	VendorID   string   `json:"vendor_id,omitempty"`
	Family     string   `json:"family,omitempty"`
	Model      string   `json:"model,omitempty"`
	Stepping   int32    `json:"stepping,omitempty"`
	PhysicalID string   `json:"physical_id,omitempty"`
	CoreID     string   `json:"core_id,omitempty"`
	CoresCount int32    `json:"cores_count,omitempty"`
	ModelName  string   `json:"model_name,omitempty"`
	Mhz        float64  `json:"mhz,omitempty"`
	CacheSize  int32    `json:"cache_size,omitempty"`
	Flags      []string `json:"flags,omitempty"`
	Microcode  string   `json:"microcode,omitempty"`
}

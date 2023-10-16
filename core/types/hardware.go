package types

type Host struct {
	// OS is the operating system name.
	// Example: "linux"
	OS string `json:"os"`

	// Arch is the hardware architecture.
	// Example: "arm64"
	Arch string `json:"arch"`

	// Platform is the platform name.
	// Example: "arch"
	Platform string `json:"platform"`

	// Version is the platform version.
	// Example: "13.5.2"
	Version string `json:"version"`

	// Name is the hostname.
	// Example: "my-host"
	Name string `json:"name"`
}

type Hardware struct {
	// Dockerized is true if the application is running inside a Docker container.
	Dockerized bool `json:"dockerized"`

	// Host is the host information.
	Host Host `json:"host"`
}

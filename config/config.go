package config

import (
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/storage"
)

var Current = New()

type Config struct {
	HostVertex string `json:"host"`
	PortVertex string `json:"port"`

	HostKernel string `json:"host_kernel"`
	PortKernel string `json:"port_kernel"`
}

func New() Config {
	return Config{
		HostVertex: "127.0.0.1:6130",
		PortVertex: "6130",

		HostKernel: "http://localhost:6131",
		PortKernel: "6131",
	}
}

func (c Config) Apply() error {
	configJsContent := fmt.Sprintf("window.apiURL = \"http://%s\";", c.HostVertex)
	return os.WriteFile(path.Join(storage.Path, "client", "dist", "config.js"), []byte(configJsContent), os.ModePerm)
}

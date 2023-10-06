package config

import (
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/storage"
)

const urlFormat = "http://%s:%s"

var Current = New()

type Config struct {
	Host string `json:"host"`

	Port           string `json:"port"`
	PortKernel     string `json:"port_kernel"`
	PortProxy      string `json:"port_proxy"`
	PortPrometheus string `json:"port_prometheus"`
}

func New() Config {
	return Config{
		Host: "127.0.0.1",

		Port:           "6130",
		PortKernel:     "6131",
		PortProxy:      "80",
		PortPrometheus: "2112",
	}
}

func (c Config) VertexURL() string {
	return fmt.Sprintf(urlFormat, c.Host, c.Port)
}

func (c Config) KernelURL() string {
	return fmt.Sprintf(urlFormat, c.Host, c.PortKernel)
}

func (c Config) ProxyURL() string {
	return fmt.Sprintf(urlFormat, c.Host, c.PortProxy)
}

func (c Config) Apply() error {
	configJsContent := fmt.Sprintf("window.apiURL = \"%s\";", c.VertexURL())
	return os.WriteFile(path.Join(storage.Path, "client", "dist", "config.js"), []byte(configJsContent), os.ModePerm)
}

package config

import (
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/storage"
)

const urlFormat = "http://%s:%s"

var Current = New()

type Mode string

const (
	ProductionMode Mode = "production"
	DebugMode      Mode = "debug"
	EmptyStr            = ""
)

type Config struct {
	mode Mode

	Host string `json:"host"`

	Port           string `json:"port"`
	PortKernel     string `json:"port_kernel"`
	PortProxy      string `json:"port_proxy"`
	PortPrometheus string `json:"port_prometheus"`
}

func New() Config {
	host, err := net.LocalIP()
	if err != nil {
		log.Error(err)
		host = "127.0.0.1"
	}

	c := Config{
		mode: ProductionMode,

		Host: host,

		Port:           "6130",
		PortKernel:     "6131",
		PortProxy:      "80",
		PortPrometheus: "2112",
	}

	if os.Getenv("DEBUG") == "1" {
		log.Warn("debug mode enabled. proceed with caution!")
		c.mode = DebugMode
	}

	return c
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

func (c Config) Debug() bool {
	return c.mode == DebugMode
}

func (c Config) Apply() error {
	configJsContent := fmt.Sprintf("window.apiURL = \"%s\";", c.VertexURL())
	return os.WriteFile(path.Join(storage.Path, "client", "dist", "config.js"), []byte(configJsContent), os.ModePerm)
}

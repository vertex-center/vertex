package config

import (
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/storage"
)

var Current = New()

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func New() Config {
	return Config{
		Host: "127.0.0.1:6130",
		Port: "6130",
	}
}

func (c Config) Apply() error {
	configJsContent := fmt.Sprintf("window.apiURL = \"http://%s\";", c.Host)
	return os.WriteFile(path.Join(storage.PathClient, "dist", "config.js"), []byte(configJsContent), os.ModePerm)
}

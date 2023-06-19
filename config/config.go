package config

var Current = New()

type Config struct {
	Host string `json:"host"`
}

func New() Config {
	return Config{
		Host: "127.0.0.1:6130",
	}
}

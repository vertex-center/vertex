package config

type Config struct {
	Port string `json:"port"`
}

func New() Config {
	return Config{
		Port: "6130",
	}
}

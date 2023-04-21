package types

type EnvVariable struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Secret      bool   `json:"secret,omitempty"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type Service struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Repository describes where to find the service.
	// - Example for GitHub: github.com/vertex-center/vertex-spotify
	// - Example for LocalStorage: localstorage:/Users/username/service
	Repository   string          `json:"repository"`
	Description  string          `json:"description"`
	EnvVariables []EnvVariable   `json:"environment,omitempty"`
	Dependencies map[string]bool `json:"dependencies,omitempty"`
}

type ServiceRepository interface {
	GetAvailable() []Service
}

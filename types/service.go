package types

const (
	URLKindClient = "client"
)

type EnvVariable struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Secret      bool   `json:"secret,omitempty"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type ServiceMethodScript struct {
	Filename     string           `json:"file"`
	Dependencies *map[string]bool `json:"dependencies,omitempty"`
}

type ServiceMethodRelease struct {
	Dependencies *map[string]bool `json:"dependencies,omitempty"`
}

type ServiceMethodDocker struct {
	Image      *string `json:"image,omitempty"`
	Dockerfile *string `json:"dockerfile,omitempty"`

	// Ports is a map containing docker port as a key, and output port as a value.
	// The output port is automatically adjusted with PORT environment variables.
	Ports   *map[string]string `json:"ports,omitempty"`
	Volumes *map[string]string `json:"volumes,omitempty"`
}

type ServiceMethods struct {
	Script  *ServiceMethodScript  `json:"script,omitempty"`
	Release *ServiceMethodRelease `json:"release,omitempty"`
	Docker  *ServiceMethodDocker  `json:"docker,omitempty"`
}

type URL struct {
	Name      string  `json:"name"`
	Port      string  `json:"port"`
	PingRoute *string `json:"ping,omitempty"`
	Kind      string  `json:"kind"`
}

type Service struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Repository describes where to find the service.
	// - Example for GitHub: github.com/vertex-center/vertex-spotify
	// - Example for LocalStorage: localstorage:/Users/username/service
	Repository     string         `json:"repository"`
	Description    string         `json:"description"`
	EnvDefinitions []EnvVariable  `json:"environment,omitempty"`
	URLs           []URL          `json:"urls,omitempty"`
	Methods        ServiceMethods `json:"methods"`
}

type ServiceRepository interface {
	Get(repo string) (Service, error)
	GetAll() []Service
}

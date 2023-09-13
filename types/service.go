package types

import "errors"

const (
	URLKindClient = "client"
)

type Service struct {
	// ID is the identifier of the service. It must be unique.
	ID string `yaml:"id" json:"id"`

	// Name is the displayed name of the service.
	Name string `yaml:"name" json:"name"`

	// Repository is the url of the repository, if it is an external repository.
	Repository *string `yaml:"repository,omitempty" json:"repository,omitempty"`

	// Description describes the service in a few words.
	Description string `yaml:"description" json:"description"`

	// Features describes some features of the service to help Vertex.
	Features *Features `yaml:"features,omitempty" json:"features,omitempty"`

	// EnvDefinitions defines all parameterizable environment variables.
	EnvDefinitions []EnvDefinition `yaml:"environment,omitempty" json:"environment,omitempty"`

	// URLs defines all service urls.
	URLs []URL `yaml:"urls,omitempty" json:"urls,omitempty"`

	// Methods defines different methods to install the service.
	Methods ServiceMethods `yaml:"methods" json:"methods"`
}

type DatabaseFeature struct {
	// The database Type. Can be redis, postgres...
	Type string `yaml:"type" json:"type"`

	// The database Port. Must be the name
	// of an environment variable.
	Port string `yaml:"port" json:"port"`

	// The Username to connect to the database. Must be the name
	// of an environment variable.
	Username *string `yaml:"username" json:"username"`

	// The Password to connect to the database. Must be the name
	// of an environment variable.
	Password *string `yaml:"password" json:"password"`
}

type Features struct {
	// The database feature describes the database made available
	// by this service.
	Databases *[]DatabaseFeature `yaml:"databases" json:"databases"`
}

type EnvDefinition struct {
	// Type is the environment variable type.
	// It can be: port, string, url.
	Type string `yaml:"type" json:"type"`

	// Name is the environment variable name that will be used by the service.
	Name string `yaml:"name" json:"name"`

	// DisplayName is a readable name for the user.
	DisplayName string `yaml:"display_name" json:"display_name"`

	// Secret is true if the value should not be read.
	Secret *bool `yaml:"secret,omitempty" json:"secret,omitempty"`

	// Default defines a default value.
	Default string `yaml:"default,omitempty" json:"default,omitempty"`

	// Description describes this variable to the user.
	Description string `yaml:"description" json:"description"`
}

type ServiceDependency struct{}

type ServiceClone struct {
	Repository string `yaml:"repository" json:"repository"`
}

type ServiceMethodScript struct {
	// Filename is the name of the file to run to start the service.
	Filename string `yaml:"file" json:"file"`

	// Clone describes the repository to clone if some files are needed to run the script.
	Clone *ServiceClone `yaml:"clone,omitempty" json:"clone,omitempty"`

	// Dependencies lists all dependencies needed before running the service.
	Dependencies *map[string]ServiceDependency `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
}

type ServiceMethodRelease struct {
	// Dependencies lists all dependencies needed before running the service.
	Dependencies *map[string]ServiceDependency `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
}

type ServiceMethodDocker struct {
	// Image is the Docker image to run.
	Image *string `yaml:"image,omitempty" json:"image,omitempty"`

	// Clone describes the repository to clone if some files are needed to run the script.
	Clone *ServiceClone `yaml:"clone,omitempty" json:"clone,omitempty"`

	// Dockerfile is the name of the Dockerfile if the repository is cloned.
	Dockerfile *string `yaml:"dockerfile,omitempty" json:"dockerfile,omitempty"`

	// Ports is a map containing docker port as a key, and output port as a value.
	// The output port is automatically adjusted with PORT environment variables.
	Ports *map[string]string `yaml:"ports,omitempty" json:"ports,omitempty"`

	// Volumes is a map containing output folder as a key, and input folder from Docker
	// as a string value.
	Volumes *map[string]string `yaml:"volumes,omitempty" json:"volumes,omitempty"`

	// Environment is a map containing docker environment variable as a key, and
	// its corresponding service environment name as a value.
	Environment *map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`

	// Capabilities is an array containing all additional Docker capabilities.
	Capabilities *[]string `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`

	// Sysctls allows to modify kernel parameters.
	Sysctls *map[string]string `yaml:"sysctls,omitempty" json:"sysctls,omitempty"`
}

type ServiceMethods struct {
	// Script is a method to launch the service with a shell script.
	Script *ServiceMethodScript `yaml:"script,omitempty" json:"script,omitempty"`

	// Release is a method to download and launch the service with
	// precompiled binaries from GitHub.
	Release *ServiceMethodRelease `yaml:"release,omitempty" json:"release,omitempty"`

	// Docker is a method to run the service with Docker.
	Docker *ServiceMethodDocker `yaml:"docker,omitempty" json:"docker,omitempty"`
}

type URL struct {
	// Name is the name displayed to the used describing this URL.
	Name string `yaml:"name" json:"name"`

	// Port is the port where this url is supposed to be.
	// Note that this port is mapped to the default value of an environment definition if possible,
	// but the port here doesn't change with the environment.
	Port string `yaml:"port" json:"port"`

	// HomeRoute allows to specify a route to change the home path.
	HomeRoute *string `yaml:"home,omitempty" json:"home,omitempty"`

	// PingRoute allows to specify a route to change the ping path.
	PingRoute *string `yaml:"ping,omitempty" json:"ping,omitempty"`

	// Kind is the type of url.
	// It can be: client, server.
	Kind string `yaml:"kind" json:"kind"`
}

var (
	ErrServiceNotFound = errors.New("the service was not found")
)

type ServiceAdapterPort interface {
	// Get a service with its id. Returns ErrServiceNotFound if
	// the service was not found.
	Get(id string) (Service, error)

	GetScript(id string) ([]byte, error)

	// GetAll gets all available services.
	GetAll() []Service

	// Reload the adapter
	Reload() error
}

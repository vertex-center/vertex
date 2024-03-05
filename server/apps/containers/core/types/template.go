package types

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

const MaxSupportedVersion Version = 3

var ErrTemplateNotFound = errors.NotFoundf("template")

type Version int

type TemplateVersioning struct {
	// Version of the template format used.
	Version Version `yaml:"version" json:"version" example:"3"`
}

type Template struct {
	TemplateVersioning `yaml:",inline"`

	// ID is the identifier of the template. It must be unique.
	ID string `yaml:"id" json:"id" example:"vertex-auth"`

	// Name is the displayed name of the template.
	Name string `yaml:"name" json:"name" example:"Vertex Auth"`

	// Repository is the url of the repository, if it is an external repository.
	Repository *string `yaml:"repository,omitempty" json:"repository,omitempty" example:"https://github.com/vertex-center/vertex"`

	// Description describes the template in a few words.
	Description string `yaml:"description" json:"description" example:"The authentication backend of Vertex."`

	// Color is the main color of the template.
	Color *string `yaml:"color,omitempty" json:"color,omitempty" example:"#f38ba8"`

	// Icon is the icon link of the template, located in ./live/templates/icons/.
	Icon *string `yaml:"icon,omitempty" json:"icon,omitempty" example:"vertex.svg"`

	// Features describes some features of the template to help Vertex.
	Features *Features `yaml:"features,omitempty" json:"features,omitempty"`

	// Env defines all parameterizable environment variables.
	Env []TemplateEnv `yaml:"environment,omitempty" json:"environment,omitempty"`

	// Databases defines all databases used by the template.
	Databases map[string]DatabaseEnvironment `yaml:"databases,omitempty" json:"databases,omitempty"`

	// Ports defines all ports that should be exposed by the container.
	Ports []TemplatePort `yaml:"ports,omitempty" json:"ports,omitempty"`

	// URLs defines all template urls.
	// Deprecated: URLs are deleted in version 3.
	URLs []URL `yaml:"urls,omitempty" json:"urls,omitempty"`

	// Methods define different methods to install the template.
	Methods TemplateMethods `yaml:"methods" json:"methods"`
}

type TemplateV1 Template

// Upgrade TemplateV1 to TemplateV2.
// Ports are now a map from port:ENV_NAME instead of port:port.
func (s *TemplateV1) Upgrade() *TemplateV2 {
	s.Version = 2
	if s.Methods.Docker != nil && s.Methods.Docker.Ports != nil {
		ports := make(map[string]string)
		for in, out := range *s.Methods.Docker.Ports {
			for _, e := range s.Env {
				if e.Type == "port" && e.Default == out {
					ports[in] = e.Name
					break
				}
			}
		}
		s.Methods.Docker.Ports = &ports
	}
	for i, url := range s.URLs {
		for _, e := range s.Env {
			if e.Type == "port" && e.Default == url.Port {
				s.URLs[i].Port = e.Name
				break
			}
		}
	}
	return (*TemplateV2)(s)
}

type TemplateV2 Template

// Upgrade TemplateV2 to TemplateV3.
// Ports are now an array in the root of the template.
func (s *TemplateV2) Upgrade() *TemplateV3 {
	s.Version = 3
	envs := make([]TemplateEnv, 0)
	for _, env := range s.Env {
		if env.Type == "port" {
			s.Ports = append(s.Ports, TemplatePort{
				Name: env.Name,
				Port: env.Default,
			})
		} else {
			envs = append(envs, env)
		}
	}
	s.Env = envs
	s.URLs = []URL{}
	return (*TemplateV3)(s)
}

type TemplateV3 Template

func (s *TemplateV3) Upgrade() *Template {
	return (*Template)(s)
}

func (s *Template) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var template struct {
		TemplateVersioning `yaml:",inline"`
		ID                 string `yaml:"id"`
	}

	err := unmarshal(&template)
	if err != nil {
		return err
	}
	s.TemplateVersioning = template.TemplateVersioning

	log.Debug("reading template",
		vlog.Int("version", int(template.Version)),
		vlog.String("id", template.ID),
	)

	var tmpl any
	switch template.Version {
	case 0, 1:
		tmpl = &TemplateV1{}
	case 2:
		tmpl = &TemplateV2{}
	case 3:
		tmpl = &TemplateV3{}
	}
	err = unmarshal(tmpl)
	if err != nil {
		return err
	}

	version := template.Version

	switch version {
	case 0, 1:
		tmpl = tmpl.(*TemplateV1).Upgrade()
		fallthrough
	case 2:
		tmpl = tmpl.(*TemplateV2).Upgrade()
		fallthrough
	case 3:
		tmpl = tmpl.(*TemplateV3).Upgrade()
	}

	if serv, ok := tmpl.(*Template); ok {
		*s = *serv
	} else {
		return fmt.Errorf("unknown template version: %d", version)
	}
	return nil
}

type TemplateUpdate struct {
	Available bool `json:"available"`
}

type DatabaseEnvironment struct {
	// DisplayName is a readable name for the user.
	DisplayName string `yaml:"display_name" json:"display_name"`

	// The database Types. Can be redis, postgres...
	Types []string `yaml:"types" json:"types"`

	// The database environment names.
	Names DatabaseEnvironmentNames `yaml:"names" json:"names"`
}

type DatabaseEnvironmentNames struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
}

type DatabaseFeature struct {
	// The database Type. Can be redis, postgres...
	Type string `yaml:"type" json:"type" example:"postgres"`

	// The database Category. Can be 'sql', 'redis'...
	Category string `yaml:"category" json:"category" example:"sql"`

	// The database Port. Must be the name
	// of an environment variable.
	Port string `yaml:"port" json:"port" example:"5432"`

	// The Username to connect to the database. Must be the name
	// of an environment variable.
	Username *string `yaml:"username" json:"username" example:"postgres"`

	// The Password to connect to the database. Must be the name
	// of an environment variable.
	Password *string `yaml:"password" json:"password" example:"postgres"`

	// The DefaultDatabase to connect to the database. Must be the name
	// of an environment variable.
	DefaultDatabase *string `yaml:"default-database" json:"database_default" example:"postgres"`
}

type Features struct {
	// The database feature describes the database made available
	// by this template.
	Databases *[]DatabaseFeature `yaml:"databases" json:"databases"`
}

type TemplateEnv struct {
	// Type is the environment variable type.
	// It can be: port, string, url.
	Type string `yaml:"type" json:"type" example:"port"`

	// Name is the environment variable name that will be used by the template.
	Name string `yaml:"name" json:"name" example:"PORT"`

	// DisplayName is a readable name for the user.
	DisplayName string `yaml:"display_name" json:"display_name" example:"Server Port"`

	// Secret is true if the value should not be read.
	Secret *bool `yaml:"secret,omitempty" json:"secret,omitempty" example:"false"`

	// Default defines a default value.
	Default string `yaml:"default,omitempty" json:"default,omitempty" example:"8080"`

	// Description describes this variable to the user.
	Description string `yaml:"description" json:"description" example:"The port where the server will listen."`
}

type TemplatePort struct {
	// Name is the name displayed to the used describing this port.
	Name string `yaml:"name" json:"name" example:"Server Port"`

	// Port is the port where this port is supposed to be.
	Port string `yaml:"port" json:"port" example:"3000"`
}

type TemplateDependency struct{}

type TemplateClone struct {
	Repository string `yaml:"repository" json:"repository" example:"https://github.com/vertex-center/vertex"`
}

type TemplateMethodDocker struct {
	// Image is the Docker image to run.
	Image *string `yaml:"image,omitempty" json:"image,omitempty" example:"ghcr.io/vertex-center/vertex"`

	// Clone describes the repository to clone if some files are needed to run the script.
	Clone *TemplateClone `yaml:"clone,omitempty" json:"clone,omitempty"`

	// Dockerfile is the name of the Dockerfile if the repository is cloned.
	Dockerfile *string `yaml:"dockerfile,omitempty" json:"dockerfile,omitempty" example:"Dockerfile"`

	// Ports is a map containing docker port as a key, and output port as a value.
	// The output port is automatically adjusted with PORT environment variables.
	// Deprecated: Use the root Ports variable instead.
	Ports *map[string]string `yaml:"ports,omitempty" json:"ports,omitempty"`

	// Volumes is a map containing output folder as a key, and input folder from Docker
	// as a string value.
	Volumes *map[string]string `yaml:"volumes,omitempty" json:"volumes,omitempty"`

	// Capabilities is an array containing all additional Docker capabilities.
	Capabilities *[]string `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`

	// Sysctls allows to modify kernel parameters.
	Sysctls *map[string]string `yaml:"sysctls,omitempty" json:"sysctls,omitempty"`

	// Cmd is the command to run in the container.
	Cmd *string `yaml:"command,omitempty" json:"command,omitempty"`
}

type TemplateMethods struct {
	// Docker is a method to run the template with Docker.
	Docker *TemplateMethodDocker `yaml:"docker,omitempty" json:"docker,omitempty"`
}

// Deprecated: Deleted in version 3.
type URL struct {
	// Name is the name displayed to the used describing this URL.
	Name string `yaml:"name" json:"name" example:"Vertex Client"`

	// Port is the port where this url is supposed to be.
	// Note that this port is mapped to the default value of an environment definition if possible,
	// but the port here doesn't change with the environment.
	Port string `yaml:"port" json:"port" example:"3000"`

	// HomeRoute allows specifying a route to change the home path.
	HomeRoute *string `yaml:"home,omitempty" json:"home,omitempty" example:"/home"`

	// PingRoute allows specifying a route to change the ping path.
	PingRoute *string `yaml:"ping,omitempty" json:"ping,omitempty" example:"/ping"`

	// Kind is the type of url.
	// It can be: client, server.
	Kind string `yaml:"kind" json:"kind" enum:"client,server"`
}

type SetDatabasesOptions struct {
	// The database name to connect to the database. Must be the name
	// of an environment variable.
	DatabaseName *string
}

package types

import "github.com/google/uuid"

type ContainerSettings struct {
	// Method indicates how the container is installed.
	// It can be by script, release or docker.
	InstallMethod *string `json:"install_method,omitempty" yaml:"install_method,omitempty"`

	// LaunchOnStartup indicates if the container needs to start automatically when Vertex starts.
	// The default value is true.
	LaunchOnStartup *bool `json:"launch_on_startup,omitempty" yaml:"launch_on_startup,omitempty"`

	// DisplayName is a custom name for the container.
	DisplayName string `json:"display_name" yaml:"display_name"`

	// Databases describe the databases used by the container.
	// The key is the database ID, and the value is the database container UUID.
	Databases map[string]uuid.UUID `json:"databases,omitempty" yaml:"databases,omitempty"`

	// Version is the version of the program.
	Version *string `json:"version,omitempty" yaml:"version,omitempty"`

	// Tags are the tags assigned to the container.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

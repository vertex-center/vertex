package dependencies

import (
	"github.com/vertex-center/vertex/dependencies/dependency"
)

func Get(id string) (*dependency.Dependency, error) {
	return dependency.New(id)
}

package adapter

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/storage"
)

const InstanceEnvPath = ".env"

type InstanceEnvFSAdapter struct {
	instancesPath string
}

type InstanceEnvFSAdapterParams struct {
	instancesPath string
}

func NewInstanceEnvFSAdapter(params *InstanceEnvFSAdapterParams) instancestypes.InstanceEnvAdapterPort {
	if params == nil {
		params = &InstanceEnvFSAdapterParams{}
	}
	if params.instancesPath == "" {
		params.instancesPath = path.Join(storage.Path, "instances")
	}

	adapter := &InstanceEnvFSAdapter{
		instancesPath: params.instancesPath,
	}

	return adapter
}

func (a *InstanceEnvFSAdapter) Save(uuid uuid.UUID, env instancestypes.InstanceEnvVariables) error {
	envPath := path.Join(a.instancesPath, uuid.String(), InstanceEnvPath)

	file, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	for key, value := range env {
		_, err := file.WriteString(strings.Join([]string{key, value}, "=") + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *InstanceEnvFSAdapter) Load(uuid uuid.UUID) (instancestypes.InstanceEnvVariables, error) {
	envPath := path.Join(a.instancesPath, uuid.String(), InstanceEnvPath)

	file, err := os.Open(envPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := instancestypes.InstanceEnvVariables{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) == 1 {
			env[line[0]] = ""
			continue
		}
		if len(line) == 2 {
			env[line[0]] = line[1]
			continue
		}
		return nil, errors.New("failed to read .env")
	}

	return env, nil
}

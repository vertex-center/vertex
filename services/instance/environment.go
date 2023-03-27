package instance

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-json"
	"github.com/vertex-center/vertex/storage"
)

type EnvVariables struct {
	Entries map[string]string
}

func NewEnvVariables() *EnvVariables {
	return &EnvVariables{Entries: map[string]string{}}
}

func (i *Instance) LoadEnvFromDisk() error {
	filepath := path.Join(storage.PathInstances, i.UUID.String(), ".env")

	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) < 2 {
			return errors.New("failed to read .env")
		}

		i.EnvVariables.Entries[line[0]] = line[1]
	}

	return nil
}

func (i *Instance) SetEnv(variables map[string]string) error {
	filepath := path.Join(storage.PathInstances, i.UUID.String(), ".env")

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	for key, value := range variables {
		_, err := file.WriteString(strings.Join([]string{key, value}, "=") + "\n")
		if err != nil {
			return err
		}
	}

	i.EnvVariables.Entries = variables

	return nil
}

func (env *EnvVariables) MarshalJSON() ([]byte, error) {
	return json.Marshal(env.Entries)
}

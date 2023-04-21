package types

import "encoding/json"

type EnvVariables struct {
	Entries map[string]string
}

func NewEnvVariables() *EnvVariables {
	return &EnvVariables{Entries: map[string]string{}}
}

func (env *EnvVariables) MarshalJSON() ([]byte, error) {
	return json.Marshal(env.Entries)
}

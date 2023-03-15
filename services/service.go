package services

type EnvVariable struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Secret      bool   `json:"secret,omitempty"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type Service struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Repository   string        `json:"repository"`
	Description  string        `json:"description"`
	EnvVariables []EnvVariable `json:"environment,omitempty"`
}

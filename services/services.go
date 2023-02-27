package services

type Service struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Repository  string `json:"repository"`
	Description string `json:"description"`
}

type InstalledService struct {
	Service
	Status string `json:"status"`
}

package types

type DockerContainerInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Platform string `json:"platform"`
}

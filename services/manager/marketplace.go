package servicesmanager

import "github.com/vertex-center/vertex/services"

var available = []services.Service{
	{
		ID:         "vertex-redis",
		Name:       "Vertex Redis",
		Repository: "github.com/vertex-center/vertex-redis",
	},
	{
		ID:         "vertex-spotify",
		Name:       "Vertex Spotify",
		Repository: "github.com/vertex-center/vertex-spotify",
	},
}

func ListAvailable() []services.Service {
	return available
}

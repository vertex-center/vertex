package servicesmanager

import "github.com/vertex-center/vertex/services"

var available = []services.Service{
	{
		ID:          "vertex-redis",
		Name:        "Vertex Redis",
		Repository:  "github.com/vertex-center/vertex-redis",
		Description: "A Redis wrapper for Vertex.",
	},
	{
		ID:          "vertex-spotify",
		Name:        "Vertex Spotify",
		Repository:  "github.com/vertex-center/vertex-spotify",
		Description: "This Spotify service collects all your spotify listening and can publish player events on Redis.",
	},
}

func ListAvailable() []services.Service {
	return available
}

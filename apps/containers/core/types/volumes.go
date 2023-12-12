package types

type (
	Volumes []Volume
	Volume  struct {
		In  string `json:"in"`  // Path in the container
		Out string `json:"out"` // Path on the host
	}
)

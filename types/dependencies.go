package types

type Dependency struct {
	*Package

	Installed bool `json:"installed"`
}

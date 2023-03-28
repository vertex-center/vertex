package dependencies

type Package struct {
	Name        string
	Description string
	Homepage    string
	License     string
	Check       string
	Install     map[string]string
}

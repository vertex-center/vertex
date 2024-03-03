//go:build darwin

package user

func getUsersDir() string {
	return "/Users"
}

func validateUsername(username string) bool {
	return username != "Shared"
}

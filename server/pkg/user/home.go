//go:build !darwin

package user

func getUsersDir() string {
	return "/home"
}

func validateUsername(username string) bool {
	return true
}

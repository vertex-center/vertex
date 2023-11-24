package user

import "os"

type User struct {
	Name    string
	HomeDir string
}

func GetAll() ([]User, error) {
	users := []User{{
		Name:    "root",
		HomeDir: "/root",
	}}

	usersDir := getUsersDir()

	homeEntries, err := os.ReadDir(usersDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range homeEntries {
		if entry.IsDir() {
			username := entry.Name()
			valid := validateUsername(username)
			if !valid {
				continue
			}
			users = append(users, User{
				Name:    username,
				HomeDir: usersDir + "/" + username,
			})
		}
	}

	return users, nil
}

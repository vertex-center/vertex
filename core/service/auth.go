package service

import (
	"encoding/base64"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	adapter port.AuthAdapter
}

func NewAuthService(adapter port.AuthAdapter) port.AuthService {
	return &AuthService{
		adapter: adapter,
	}
}

func (s *AuthService) Login(login, password string) error {
	return nil
}

func (s *AuthService) Register(login, password string) error {
	// TODO: make these settings configurable in the admin settings
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id

	it := uint32(3)
	mem := uint32(12 * 1024)
	threads := uint8(4)
	salt := "vertex"

	key := argon2.IDKey([]byte(password), []byte(salt), it, mem, threads, 32)
	hash := base64.StdEncoding.EncodeToString(key)

	return s.adapter.CreateAccount(login, types.CredentialsArgon2id{
		Login:       login,
		Hash:        hash,
		Type:        "argon2id",
		Iterations:  it,
		Memory:      mem,
		Parallelism: threads,
		Salt:        salt,
	})
}

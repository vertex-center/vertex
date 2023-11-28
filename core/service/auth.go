package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
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

// Login checks the provided login and password against the database. If the login is
// valid and the password matches, it returns a token that can be used to authenticate
// future requests.
func (s *AuthService) Login(login, password string) (types.Token, error) {
	creds, err := s.adapter.GetCredentials(login)
	if err != nil {
		return types.Token{}, err
	}

	for _, cred := range creds {
		storedKey, err := base64.StdEncoding.DecodeString(cred.Hash)
		if err != nil {
			log.Error(errors.New("failed to decode stored key"), vlog.String("reason", err.Error()))
			continue
		}

		key := argon2.IDKey([]byte(password), []byte(cred.Salt), cred.Iterations, cred.Memory, cred.Parallelism, cred.KeyLen)
		if bytes.Equal(storedKey, key) {
			token, err := s.generateToken()
			if err != nil {
				return types.Token{}, err
			}
			token.Username = login
			err = s.adapter.SaveToken(&token)
			if err != nil {
				return types.Token{}, err
			}
			return token, nil
		}
	}

	return types.Token{}, types.ErrLoginFailed
}

// Register creates a new user account. It can return ErrLoginEmpty, ErrPasswordEmpty, or
// ErrPasswordLength if the login or password is too short.
func (s *AuthService) Register(login, password string) (types.Token, error) {
	// TODO: make these settings configurable in the admin settings
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id

	if len(login) == 0 {
		return types.Token{}, types.ErrLoginEmpty
	}
	if len(password) == 0 {
		return types.Token{}, types.ErrPasswordEmpty
	}
	if len(password) < 8 {
		return types.Token{}, types.ErrPasswordLength
	}

	it := uint32(3)
	mem := uint32(12 * 1024)
	threads := uint8(4)
	salt := "vertex"
	keyLen := uint32(32)

	key := argon2.IDKey([]byte(password), []byte(salt), it, mem, threads, keyLen)
	hash := base64.StdEncoding.EncodeToString(key)

	cred := types.CredentialsArgon2id{
		Login:       login,
		Hash:        hash,
		Type:        "argon2id",
		Iterations:  it,
		Memory:      mem,
		Parallelism: threads,
		Salt:        salt,
		KeyLen:      keyLen,
	}
	err := s.adapter.CreateAccount(login, cred)
	if err != nil {
		return types.Token{}, err
	}

	return s.Login(login, password)
}

func (s *AuthService) Logout(token string) error {
	return s.adapter.RemoveToken(token)
}

func (s *AuthService) Verify(token string) error {
	_, err := s.adapter.GetToken(token)
	if err != nil {
		return types.ErrTokenInvalid
	}
	return nil
}

func (s *AuthService) generateToken() (types.Token, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return types.Token{}, err
	}
	return types.Token{
		Token: base64.StdEncoding.EncodeToString(token),
	}, nil
}

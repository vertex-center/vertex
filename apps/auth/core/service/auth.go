package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/config"
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
func (s *AuthService) Login(login, password string) (types.Session, error) {
	creds, err := s.adapter.GetCredentials(login)
	if err != nil {
		return types.Session{}, err
	}

	for _, cred := range creds {
		storedKey, err := base64.StdEncoding.DecodeString(cred.Hash)
		if err != nil {
			log.Error(errors.New("failed to decode stored key"), vlog.String("reason", err.Error()))
			continue
		}

		key := argon2.IDKey([]byte(password), []byte(cred.Salt), cred.Iterations, cred.Memory, cred.Parallelism, cred.KeyLen)
		if bytes.Equal(storedKey, key) {
			session, err := s.generateToken()
			if err != nil {
				return types.Session{}, err
			}

			users, err := s.adapter.GetUsersByCredential(cred.ID)
			if err != nil {
				return types.Session{}, err
			}
			if len(users) == 0 {
				return types.Session{}, errors.New("no users linked to this credential")
			}
			session.UserID = users[0].ID
			err = s.adapter.SaveSession(&session)
			if err != nil {
				return types.Session{}, err
			}
			return session, nil
		}
	}

	return types.Session{}, types.ErrLoginFailed
}

// Register creates a new user account. It can return ErrLoginEmpty, ErrPasswordEmpty, or
// ErrPasswordLength if the login or password is too short.
func (s *AuthService) Register(login, password string) (types.Session, error) {
	// TODO: make these settings configurable in the admin settings
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id

	if len(login) == 0 {
		return types.Session{}, types.ErrLoginEmpty
	}
	if len(password) == 0 {
		return types.Session{}, types.ErrPasswordEmpty
	}
	if len(password) < 8 {
		return types.Session{}, types.ErrPasswordLength
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
		return types.Session{}, err
	}

	return s.Login(login, password)
}

func (s *AuthService) Logout(token string) error {
	return s.adapter.DeleteSession(token)
}

func (s *AuthService) Verify(token string) (*types.Session, error) {
	if token == config.Current.MasterApiKey {
		log.Debug("master key used for authentication")
		return &types.Session{
			Token: token,
		}, nil
	}
	session, err := s.adapter.GetSession(token)
	if err != nil {
		return nil, types.ErrTokenInvalid
	}
	return session, nil
}

func (s *AuthService) generateToken() (types.Session, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return types.Session{}, err
	}
	return types.Session{
		Token: base64.StdEncoding.EncodeToString(token),
	}, nil
}

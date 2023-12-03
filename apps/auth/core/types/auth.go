package types

import (
	"errors"
)

// User and CredentialsArgon2id are many-to-many relationship, because one user
// can have multiple login methods, and one login method can be used to
// connect to multiple users at once.

var (
	ErrLoginEmpty         = errors.New("login is empty")
	ErrPasswordEmpty      = errors.New("password is empty")
	ErrPasswordLength     = errors.New("password length requirement not met")
	ErrLoginFailed        = errors.New("login failed")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrEmailEmpty         = errors.New("email is empty")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type CredentialsType string

const CredentialsTypeLoginPassword CredentialsType = "Login and password"

type (
	User struct {
		ID        uint   `json:"id" db:"id"`
		Username  string `json:"username" db:"username"`
		CreatedAt int64  `json:"created_at" db:"created_at"`
		UpdatedAt int64  `json:"updated_at" db:"updated_at"`
		DeletedAt *int64 `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	Email struct {
		ID        uint   `json:"id" db:"id"`
		UserID    uint   `json:"user_id" db:"user_id"`
		Email     string `json:"email" db:"email"`
		CreatedAt int64  `json:"created_at" db:"created_at"`
		UpdatedAt int64  `json:"updated_at" db:"updated_at"`
		DeletedAt *int64 `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	CredentialsMethods struct {
		Name        CredentialsType `json:"name"`
		Description string          `json:"description"`
	}

	// CredentialsArgon2id is the login method that uses Argon2id algorithm to
	// hash the password.
	CredentialsArgon2id struct {
		ID          uint   `json:"id" db:"id"`
		Login       string `json:"login" db:"login"`
		Hash        string `json:"hash" db:"hash"`
		Type        string `json:"type" db:"type"`
		Iterations  uint32 `json:"iterations" db:"iterations"`
		Memory      uint32 `json:"memory" db:"memory"`
		Parallelism uint8  `json:"parallelism" db:"parallelism"`
		Salt        string `json:"salt" db:"salt"`
		KeyLen      uint32 `json:"key_len" db:"key_len"`
		CreatedAt   int64  `json:"created_at" db:"created_at"`
		UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
		DeletedAt   *int64 `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	Session struct {
		ID        uint   `json:"id" db:"id"`
		Token     string `json:"token" db:"token"`
		UserID    uint   `json:"user,omitempty" db:"user_id"`
		CreatedAt int64  `json:"created_at" db:"created_at"`
		UpdatedAt int64  `json:"updated_at" db:"updated_at"`
		DeletedAt *int64 `json:"deleted_at,omitempty" db:"deleted_at"`
	}
)

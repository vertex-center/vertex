package types

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User and CredentialsArgon2id are many-to-many relationship, because one user
// can have multiple login methods, and one login method can be used to
// connect to multiple users at once.

var (
	ErrLoginEmpty     = errors.New("login is empty")
	ErrPasswordEmpty  = errors.New("password is empty")
	ErrPasswordLength = errors.New("password length requirement not met")
)

type User struct {
	Username            string                 `json:"username" gorm:"primaryKey"`
	CredentialsArgon2id []*CredentialsArgon2id `json:"credentials_argon2id" gorm:"many2many:user_credentials_argon2id;"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	DeletedAt           gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// CredentialsArgon2id is the login method that uses Argon2id algorithm to
// hash the password.
type CredentialsArgon2id struct {
	Login       string         `json:"login" gorm:"primaryKey"`
	Hash        string         `json:"hash" gorm:"not null"`
	Type        string         `json:"type" gorm:"default:argon2id;not null"`
	Iterations  uint32         `json:"iterations" gorm:"not null"`
	Memory      uint32         `json:"memory" gorm:"not null"`
	Parallelism uint8          `json:"parallelism" gorm:"not null"`
	Salt        string         `json:"salt" gorm:"not null"`
	KeyLen      uint32         `json:"key_len" gorm:"not null"`
	Users       []*User        `json:"users,omitempty" gorm:"many2many:user_credentials_argon2id;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

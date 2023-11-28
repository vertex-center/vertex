package types

import (
	"time"

	"gorm.io/gorm"
)

// User and CredentialsArgon2id are many-to-many relationship, because one user
// can have multiple login methods, and one login method can be used to
// connect to multiple users at once.

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
	Hash        string         `json:"hash"`
	Type        string         `json:"type" gorm:"default:argon2id"`
	Iterations  uint32         `json:"iterations"`
	Memory      uint32         `json:"memory"`
	Parallelism uint8          `json:"parallelism"`
	Salt        string         `json:"salt"`
	Users       []*User        `json:"users,omitempty" gorm:"many2many:user_credentials_argon2id;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

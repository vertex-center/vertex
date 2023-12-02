package adapter

import (
	"fmt"
	"time"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
)

type AuthDbAdapter struct {
	db *vtypes.DB
}

func NewAuthDbAdapter(db *vtypes.DB) port.AuthAdapter {
	return &AuthDbAdapter{
		db: db,
	}
}

func (a *AuthDbAdapter) CreateAccount(username string, credential types.CredentialsArgon2id) error {
	user := types.User{
		Username:  username,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	query, args, err := tx.BindNamed(`
		INSERT INTO users (username, created_at, updated_at)
		VALUES (:username, :created_at, :updated_at)
		RETURNING id
	`, user)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = tx.Get(&user, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	credential.CreatedAt = time.Now().Unix()
	credential.UpdatedAt = time.Now().Unix()

	query, args, err = tx.BindNamed(`
		INSERT INTO credentials_argon2 (login, hash, type, iterations, memory, parallelism, salt, key_len, created_at, updated_at)
		VALUES (:login, :hash, :type, :iterations, :memory, :parallelism, :salt, :key_len, :created_at, :updated_at)
		RETURNING id
	`, credential)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Get(&credential, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO credentials_argon2_users (credential_id, user_id)
		VALUES ($1, $2)
	`, credential.ID, user.ID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	return tx.Commit()
}

func (a *AuthDbAdapter) GetCredentials(login string) ([]types.CredentialsArgon2id, error) {
	var creds []types.CredentialsArgon2id

	rows, err := a.db.Queryx(`
		SELECT credentials_argon2.id, credentials_argon2.login, credentials_argon2.hash, credentials_argon2.type, credentials_argon2.iterations, credentials_argon2.memory, credentials_argon2.parallelism, credentials_argon2.salt, credentials_argon2.key_len, credentials_argon2.created_at, credentials_argon2.updated_at, credentials_argon2.deleted_at
		FROM credentials_argon2
		WHERE credentials_argon2.login = $1
	`, login)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cred types.CredentialsArgon2id
		err := rows.StructScan(&cred)
		if err != nil {
			return nil, err
		}
		creds = append(creds, cred)
	}

	return creds, nil
}

func (a *AuthDbAdapter) GetUsersByCredential(credentialID uint) ([]types.User, error) {
	var users []types.User
	err := a.db.Select(&users, `
		SELECT users.id, users.username, users.created_at, users.updated_at, users.deleted_at
		FROM users
		JOIN credentials_argon2_users ON users.id = credentials_argon2_users.user_id
		WHERE credentials_argon2_users.credential_id = $1
	`, credentialID)
	return users, err
}

func (a *AuthDbAdapter) SaveSession(session *types.Session) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	session.CreatedAt = time.Now().Unix()
	session.UpdatedAt = time.Now().Unix()

	_, err = tx.NamedExec(`
		INSERT INTO sessions (token, user_id, created_at, updated_at)
		VALUES (:token, :user_id, :created_at, :updated_at)
	`, session)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	*session = types.Session{Token: session.Token}

	err = tx.Get(session, `
		SELECT sessions.id, sessions.token, users.id, sessions.created_at, sessions.updated_at, sessions.deleted_at
		FROM sessions
		JOIN users ON sessions.user_id = users.id
		WHERE sessions.deleted_at IS NULL AND sessions.token = $1
	`, session.Token)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (a *AuthDbAdapter) DeleteSession(token string) error {
	_, err := a.db.Exec("UPDATE sessions SET deleted_at = $1 WHERE token = $2", time.Now().Unix(), token)
	return err
}

func (a *AuthDbAdapter) GetSession(token string) (*types.Session, error) {
	var session types.Session
	err := a.db.Get(&session, `
		SELECT sessions.id, sessions.token, sessions.user_id, sessions.created_at, sessions.updated_at, sessions.deleted_at
		FROM sessions
		WHERE sessions.deleted_at IS NULL AND sessions.token = $1
	`, token)
	return &session, err
}

func (a *AuthDbAdapter) GetUser(username string) (types.User, error) {
	var user types.User
	err := a.db.Get(&user, `
		SELECT users.id, users.username, users.created_at, users.updated_at, users.deleted_at
		FROM users
		WHERE users.username = $1
	`, username)
	return user, err
}

func (a *AuthDbAdapter) GetUserByID(id uint) (types.User, error) {
	var user types.User
	err := a.db.Get(&user, `
		SELECT users.id, users.username, users.created_at, users.updated_at, users.deleted_at
		FROM users
		WHERE users.id = $1
	`, id)
	return user, err
}

func (a *AuthDbAdapter) PatchUser(user types.User) (types.User, error) {
	user.UpdatedAt = time.Now().Unix()

	tx, err := a.db.Beginx()
	if err != nil {
		return user, err
	}

	_, err = tx.NamedExec(`
		UPDATE users
		SET username = :username, updated_at = :updated_at
		WHERE id = :id
	`, user)
	if err != nil {
		_ = tx.Rollback()
		return user, err
	}

	err = tx.Get(&user, `
		SELECT users.id, users.username, users.created_at, users.updated_at, users.deleted_at
		FROM users
		WHERE users.id = $1
	`, user.ID)
	if err != nil {
		_ = tx.Rollback()
		return user, err
	}

	return user, tx.Commit()
}

func (a *AuthDbAdapter) GetUserCredentialsMethods(userID uint) ([]types.CredentialsMethods, error) {
	var credsArgon2id []types.CredentialsArgon2id
	err := a.db.Select(&credsArgon2id, `
		SELECT credentials_argon2.id, credentials_argon2.login, credentials_argon2.hash, credentials_argon2.type, credentials_argon2.iterations, credentials_argon2.memory, credentials_argon2.parallelism, credentials_argon2.salt, credentials_argon2.key_len, credentials_argon2.created_at, credentials_argon2.updated_at, credentials_argon2.deleted_at
		FROM credentials_argon2
		JOIN credentials_argon2_users ON credentials_argon2.id = credentials_argon2_users.credential_id
		WHERE credentials_argon2_users.user_id = $1
	`, userID)

	var creds []types.CredentialsMethods
	for _, cr := range credsArgon2id {
		creds = append(creds, types.CredentialsMethods{
			Name:        types.CredentialsTypeLoginPassword,
			Description: fmt.Sprintf("Your password is hashed with the %s algorithm", cr.Type),
		})
	}

	return creds, err
}

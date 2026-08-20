//this file is for register users.

package api

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var ErrUserExists = errors.New("user already exists")

// UserStore stores registered users in SQLite.
type UserStore struct {
	db *sql.DB
}

// NewUserStore opens the database and creates the users table.
func NewUserStore(path string) (*UserStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &UserStore{db: db}, nil
}

// CreateUser hashes the password and stores the user.
func (s *UserStore) CreateUser(username, password string) error {
	var exists int

	err := s.db.QueryRow(
		"SELECT 1 FROM users WHERE username = ?",
		username,
	).Scan(&exists)

	if err == nil {
		return ErrUserExists
	}

	if err != sql.ErrNoRows {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO users(username, password_hash) VALUES (?, ?)",
		username,
		string(hash),
	)

	return err
}

// Authenticate checks username + password.
func (s *UserStore) Authenticate(username, password string) bool {
	var hash string

	err := s.db.QueryRow(
		"SELECT password_hash FROM users WHERE username = ?",
		username,
	).Scan(&hash)

	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}

package models

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
	"uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(ctx context.Context, name string) (User, error) {
	stmt := `INSERT INTO users (id, created_at, updated_at, name, api_key) VALUES ($1, $2, $3, $4, $5)`

	id := uuid.New()
	now := time.Now().UTC()

	apiKey, err := GenerateAPIKey(32)
	if err != nil {
		return User{}, err
	}

	_, err = m.DB.ExecContext(ctx, stmt, id, now, now, name, apiKey)
	if err != nil {
		return User{}, err
	}

	user := User{
		Id:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		ApiKey:    apiKey,
	}

	return user, nil
}

func (m *UserModel) Get(ctx context.Context, apiKey string) (User, error) {
	query := `SELECT id, created_at, updated_at, name, api_key FROM users WHERE api_key = $1`

	row := m.DB.QueryRowContext(ctx, query, apiKey)

	var user User
	err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.ApiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return user, nil
}

func GenerateAPIKey(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

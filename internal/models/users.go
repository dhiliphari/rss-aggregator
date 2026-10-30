package models

import (
	"context"
	"database/sql"
	"time"
	"uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) CreateUser(ctx context.Context, name string) (User, error) {
	stmt := `INSERT INTO users (id, created_at, updated_at, name) VALUES ($1, $2, $3, $4)`

	id := uuid.New()
	now := time.Now().UTC()

	_, err := m.DB.ExecContext(ctx, stmt, id, now, now, name)
	if err != nil {
		return User{}, err
	}

	user := User{
		Id:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}

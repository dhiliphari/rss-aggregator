package models

import (
	"context"
	"database/sql"
	"time"
	"uuid"
)

type Feed struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}

type FeedModel struct {
	DB *sql.DB
}

func (m *FeedModel) Create(ctx context.Context, name, url string, userId uuid.UUID) (Feed, error) {
	stmt := `INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
	VALUES ($1, $2, $3, $4, $5, $6)`

	id := uuid.New()
	now := time.Now().UTC()

	_, err := m.DB.ExecContext(ctx, stmt, id, now, now, name, url, userId)
	if err != nil {
		return Feed{}, err
	}

	feed := Feed{
		Id:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		URL:       url,
		UserId:    userId,
	}

	return feed, nil
}

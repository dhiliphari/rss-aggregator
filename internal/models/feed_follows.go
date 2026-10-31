package models

import (
	"context"
	"database/sql"
	"time"
	"uuid"
)

type FeedFollow struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	FeedId    uuid.UUID `json:"feed_id"`
}

type FeedFollowModel struct {
	DB *sql.DB
}

func (m *FeedFollowModel) Create(ctx context.Context, userId, feedId uuid.UUID) (FeedFollow, error) {
	stmt := `INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
	VALUES ($1, $2, $3, $4, $5)`

	id := uuid.New()
	now := time.Now().UTC()

	_, err := m.DB.ExecContext(ctx, stmt, id, now, now, userId, feedId)
	if err != nil {
		return FeedFollow{}, err
	}

	feedFollow := FeedFollow{
		Id:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserId:    userId,
		FeedId:    feedId,
	}

	return feedFollow, nil
}

func (m *FeedFollowModel) ListByUser(ctx context.Context, userId uuid.UUID) ([]FeedFollow, error) {
	query := `SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows WHERE user_id = $1;`

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fs []FeedFollow

	for rows.Next() {
		var ff FeedFollow
		err = rows.Scan(&ff.Id, &ff.CreatedAt, &ff.UpdatedAt, &ff.UserId, &ff.FeedId)
		if err != nil {
			return nil, err
		}
		fs = append(fs, ff)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fs, nil
}

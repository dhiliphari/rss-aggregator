package models

import (
	"context"
	"database/sql"
	"errors"
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

func (m *FeedModel) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	stmt := `UPDATE feeds SET last_fetched_at = NOW(), updated_at = NOW() WHERE id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
	}

	return err
}

func (m *FeedModel) GetNextFeedToFetch(ctx context.Context, count int) ([]Feed, error) {
	query := `SELECT id, created_at, updated_at, name, url, user_id FROM feeds
	ORDER BY last_fetched_at ASC NULLS FIRST LIMIT $1`
	return m.queryFeeds(ctx, query, count)
}

func (m *FeedModel) List(ctx context.Context) ([]Feed, error) {
	query := `SELECT id, created_at, updated_at, name, url, user_id FROM feeds`
	return m.queryFeeds(ctx, query)
}

func (m *FeedModel) ListByUser(ctx context.Context, userId uuid.UUID) ([]Feed, error) {
	query := `SELECT id, created_at, updated_at, name, url, user_id FROM feeds WHERE user_id = $1`
	return m.queryFeeds(ctx, query, userId)
}

func (m *FeedModel) queryFeeds(ctx context.Context, query string, args ...any) ([]Feed, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed

	for rows.Next() {
		var feed Feed

		err = rows.Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.URL, &feed.UserId)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

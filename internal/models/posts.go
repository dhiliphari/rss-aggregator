package models

import (
	"context"
	"database/sql"
	"time"
	"uuid"
)

type Post struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	URL         string    `json:"url"`
	FeedId      uuid.UUID `json:"feed_id"`
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Create(ctx context.Context, title, description, url string, publishedAt time.Time, feedId uuid.UUID) (Post, error) {
	stmt := `
		INSERT INTO posts (
			id,
			created_at,
			updated_at,
			title,
			description,
			published_at,
			url,
			feed_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	id := uuid.New()
	now := time.Now()

	_, err := m.DB.ExecContext(ctx, stmt, id, now, now, title, description, publishedAt, url, feedId)
	if err != nil {
		return Post{}, err
	}

	post := Post{
		Id:          id,
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       title,
		Description: description,
		PublishedAt: publishedAt,
		URL:         url,
		FeedId:      feedId,
	}

	return post, nil
}

func (m *PostModel) ListByUser(ctx context.Context, userId uuid.UUID, count int) ([]Post, error) {
	query := `
		SELECT 
			posts.id, 
			posts.created_at, 
			posts.updated_at, 
			posts.title, 
			posts.description, 
			posts.published_at, 
			posts.url, 
			posts.feed_id 
		FROM posts
		JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
		WHERE feed_follows.user_id = $1
		ORDER BY posts.published_at DESC
		LIMIT $2;`

	rows, err := m.DB.QueryContext(ctx, query, userId, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.Id, &p.CreatedAt, &p.UpdatedAt, &p.Title, &p.Description, &p.PublishedAt, &p.URL, &p.FeedId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

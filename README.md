# RSS Aggregator

A small Go HTTP API for following RSS feeds. Users register, add feeds, follow feeds, and read the latest posts from the feeds they follow. A background scraper fetches 10 feeds every minute, picking the ones that were fetched longest ago first, and stores new posts in PostgreSQL.

## Requirements

- Go 1.27+
- PostgreSQL. A `docker-compose.yml` is included.
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI for migrations

## Getting started

```sh
# Start Postgres
docker compose up -d

# DSN matching docker-compose.yml
export DSN="postgres://root:pgroot@localhost:5432/rssagg?sslmode=disable"

# Apply migrations
make migrate-up DSN=$DSN

# Run the server (default addr :5000)
make run DSN=$DSN
```

Other targets: `make build` (outputs `./bin/rss-aggregator`) and `make migrate-down`.

The server accepts the flags `-addr` (default `:5000`) and `-dsn`.

## Authentication

`POST /users` returns an `api_key`. Send it on authenticated requests:

```
Authorization: ApiKey <api_key>
```

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/healthz` | | Health check |
| POST | `/users` | | Create a user: `{"name": "..."}` |
| GET | `/me` | ✓ | Get the current user |
| GET | `/feeds` | | List all feeds |
| POST | `/me/feeds` | ✓ | Create a feed: `{"name": "...", "url": "..."}` |
| GET | `/me/feeds` | ✓ | List feeds created by the current user |
| GET | `/me/feeds/following` | ✓ | List the current user's feed follows |
| POST | `/me/feeds/following` | ✓ | Follow a feed: `{"feed_id": "<uuid>"}` |
| DELETE | `/me/feeds/following/{feedFollowID}` | ✓ | Unfollow a feed |
| GET | `/me/posts?count=N` | ✓ | Latest posts from followed feeds (default 10) |

## Project layout

```
cmd/api/          HTTP server, routes, handlers, middleware, feed scraper
internal/models/  Database access (users, feeds, feed_follows, posts)
internal/rss/     RSS fetching and parsing
internal/validator/ Request validation helpers
migrations/       SQL migrations (golang-migrate)
```

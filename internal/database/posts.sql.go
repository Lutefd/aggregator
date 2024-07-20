// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, title, description, url, published_at, created_at, updated_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, title, description, url, published_at, created_at, updated_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	Title       string
	Description sql.NullString
	Url         string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Url,
		arg.PublishedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
	)
	return i, err
}

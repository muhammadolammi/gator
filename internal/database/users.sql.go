// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users ( created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, created_at, updated_at, name
`

type CreateUserParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.CreatedAt, arg.UpdatedAt, arg.Name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, updated_at, name FROM users WHERE id  = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const getUserWithName = `-- name: GetUserWithName :one
SELECT id, created_at, updated_at, name FROM users WHERE name  = $1
`

func (q *Queries) GetUserWithName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithName, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, created_at, updated_at, name FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reset = `-- name: Reset :exec
DELETE FROM users
`

func (q *Queries) Reset(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, reset)
	return err
}

const userExists = `-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE name = $1
) AS exists
`

func (q *Queries) UserExists(ctx context.Context, name string) (bool, error) {
	row := q.db.QueryRowContext(ctx, userExists, name)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

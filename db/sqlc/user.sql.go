// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username, hashed_password, full_name, email
) VALUES (
    $1, $2, $3, $4
) RETURNING id, username, hashed_password, full_name, email, password_last_changed_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordLastChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getAUser = `-- name: GetAUser :one
SELECT id, username, hashed_password, full_name, email, password_last_changed_at, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getAUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordLastChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

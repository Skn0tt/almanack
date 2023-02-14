// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: option.sql

package db

import (
	"context"
)

const getOption = `-- name: GetOption :one
SELECT
  "value"
FROM
  "option"
WHERE
  key = $1
`

func (q *Queries) GetOption(ctx context.Context, key string) (string, error) {
	row := q.db.QueryRow(ctx, getOption, key)
	var value string
	err := row.Scan(&value)
	return value, err
}

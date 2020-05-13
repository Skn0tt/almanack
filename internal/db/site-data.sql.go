// Code generated by sqlc. DO NOT EDIT.
// source: site-data.sql

package db

import (
	"context"
	"encoding/json"
)

const getSiteData = `-- name: GetSiteData :one
SELECT
  "data"
FROM
  site_data
WHERE
  "key" = $1
`

func (q *Queries) GetSiteData(ctx context.Context, key string) (json.RawMessage, error) {
	row := q.db.QueryRowContext(ctx, getSiteData, key)
	var data json.RawMessage
	err := row.Scan(&data)
	return data, err
}

const setSiteData = `-- name: SetSiteData :exec
UPDATE
  site_data
SET
  "data" = $2
WHERE
  "key" = $1
`

type SetSiteDataParams struct {
	Key  string          `json:"key"`
	Data json.RawMessage `json:"data"`
}

func (q *Queries) SetSiteData(ctx context.Context, arg SetSiteDataParams) error {
	_, err := q.db.ExecContext(ctx, setSiteData, arg.Key, arg.Data)
	return err
}
